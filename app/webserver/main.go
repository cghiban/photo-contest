package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"photo-contest/app/webserver/handlers"
	"photo-contest/business/web"
	"photo-contest/foundation/database"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	log := log.New(os.Stdout, "PHOTOC : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	var cfg struct {
		conf.Version
		Web struct {
			BindAddress  string        `conf:"default:0.0.0.0:8080"`
			SessionKey   string        `conf:"default:abc123XYZ"`
			CsrfKey      string        `conf:"default:abcqwertxyz"`
			CsrfSecure   bool          `conf:"default:false"`
			IdleTimeout  time.Duration `conf:"default:5s"`
			ReadTimeout  time.Duration `conf:"default:5s"`
			WriteTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			Path        string `conf:"default:var/db.db"`
			Mode        string `conf:"default:rw"`
			JournalMode string `conf:"default:WAL"`
			Cache       string `conf:"default:shared"`
		}
	}
	cfg.Version.SVN = build
	cfg.Version.Desc = "CSHL/DNALC"

	if err := conf.Parse(os.Args[1:], "PHOTOC", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("PHOTOC", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("PHOTOC", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting

	expvar.NewString("build").Set(build)
	log.Printf("main: Started: Application initializing: version %q", build)
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config:\n%v\n", out)

	//----------------------------------------------------
	// init DB
	db, err := database.Open(database.Config{
		Path:        cfg.DB.Path,
		Mode:        cfg.DB.Mode,
		Cache:       cfg.DB.Cache,
		JournalMode: cfg.DB.JournalMode,
	})

	if err != nil {
		//return err
		return errors.Wrap(err, "connecting to db")
	}
	defer func() {
		log.Printf("main: Database Stopping: %s", cfg.DB.Path)
		db.Close()
	}()

	//dataStore := &data.DataStore{DB: db, L: log}

	sqliteVersion, _ := database.GetSQLiteVersion(db)
	log.Println("using SQLite version", sqliteVersion)

	log.Println("about to start server on ", cfg.Web.BindAddress)

	service := handlers.NewService(log, db, cfg.Web.SessionKey)

	// auth midleware...
	authMw := handlers.NewAuth(service)

	sm := mux.NewRouter()
	sm.NotFoundHandler = web.WrapMiddleware(service.NotFoundHandler, authMw.UserViaSession)
	sm.Handle("/", web.WrapMiddleware(service.Index, authMw.UserViaSession))
	sm.Handle("/guidelines", web.WrapMiddleware(service.About, authMw.UserViaSession))

	sm.Handle("/settings", web.WrapMiddleware(service.Settings, authMw.UserViaSession, authMw.RequireUser))
	//sm.Handle("/updategroup/{id:[0-9]+}", web.WrapMiddleware(service.UpdateGroup, authMw.UserViaSession, authMw.RequireUser)).Methods("POST").HeadersRegexp("Content-Type", "application/json")

	// TODO make sure we set Secure to true for production
	csrfMiddleware := csrf.Protect([]byte(cfg.Web.CsrfKey), csrf.Secure(cfg.Web.CsrfSecure))
	userRouter := sm.Methods("POST", "GET").Subrouter()
	userRouter.Use(csrfMiddleware)
	userRouter.HandleFunc("/register", service.UserSignUp)
	userRouter.HandleFunc("/login", service.UserLogIn)
	userRouter.HandleFunc("/logout", service.UserLogOut)
	userRouter.HandleFunc("/forgotpass", service.UserForgotPassword)
	userRouter.HandleFunc("/resetpass", service.UserResetPassword)
	userRouter.Handle("/profile", web.WrapMiddleware(service.UserUpdateProfile, authMw.UserViaSession, authMw.RequireUser))
	userRouter.Handle("/password", web.WrapMiddleware(service.UserUpdatePassword, authMw.UserViaSession, authMw.RequireUser))
	userRouter.Handle("/entry", web.WrapMiddleware(service.UserPhotoUpload, authMw.UserViaSession, authMw.RequireUser))
	userRouter.Handle("/submissions", web.WrapMiddleware(service.UserPhotos, authMw.UserViaSession, authMw.RequireUser))
	userRouter.Handle("/gallery", web.WrapMiddleware(service.ContestPhotos, authMw.UserViaSession, authMw.RequireUser, authMw.RequireAdmin))
	userRouter.Handle("/update_entry_status", web.WrapMiddleware(service.ContestEntryUpdateStatus, authMw.UserViaSession, authMw.RequireUser, authMw.RequireAdmin))
	userRouter.Handle("/withdraw_entry", web.WrapMiddleware(service.UserWithdrawPhoto, authMw.UserViaSession, authMw.RequireUser))
	sm.PathPrefix("/tmp/").Handler(http.StripPrefix("/tmp/", http.FileServer(http.Dir("tmp/"))))

	sm.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("var/static/"))))
	sm.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	sm.Handle("/favicon.ico", http.NotFoundHandler())

	s := &http.Server{
		Addr:         cfg.Web.BindAddress,
		Handler:      sm,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Println("Received terminate, graceful shutdown", sig)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

	return nil
}
