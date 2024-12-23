package handlers

import (
	"log"
	"net/http"
	"photo-contest/business/data/user"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

// Service data struct
type Service struct {
	log     *log.Logger
	db      *sqlx.DB
	session *sessions.CookieStore
	t       *template.Template
	//session *sqlitestore.SqliteStore
}

// NewService initializes a new Serivice
func NewService(l *log.Logger, db *sqlx.DB, sessionKey string) *Service {
	// init template
	funcMap := template.FuncMap{
		"dayToDate": func(s string) string {
			t, err := time.Parse("2006-01-02", s)
			if err != nil {
				return ""
			}
			return t.Format("Jan 2, 2006")
		},
		"dateISOish": func(t time.Time) string { return t.Format("2006-01-02 3:04pm") },
	}
	templates := template.Must(template.New("tmpls").Funcs(funcMap).ParseGlob("var/templates/*.gohtml"))
	//templates = templates.Funcs(funcMap)

	sessStore := sessions.NewCookieStore([]byte(sessionKey))
	/*sessStore, err := sqlitestore.NewSqliteStoreFromConnection(store.DB, "sessions", "/", 86400, []byte(*sessionKey))
	if err != nil {
		panic(err)
	}*/

	sessStore.Options = &sessions.Options{
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7 * 86400,
	}

	return &Service{log: l, db: db, t: templates, session: sessStore}
}

// Index - about this site
func (s *Service) Index(rw http.ResponseWriter, r *http.Request) {
	var usr *user.AuthUser
	userV := r.Context().Value("user")
	if userV != nil {
		usr = userV.(*user.AuthUser)
	}
	data := struct {
		User    *user.AuthUser
		Message string
	}{
		User:    usr,
		Message: "",
	}
	if err := s.t.ExecuteTemplate(rw, "index.gohtml", data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// About - about this site
func (s *Service) About(rw http.ResponseWriter, r *http.Request) {
	var usr *user.AuthUser
	userV := r.Context().Value("user")
	if userV != nil {
		usr = userV.(*user.AuthUser)
	}
	data := struct {
		User *user.AuthUser
	}{
		User: usr,
	}
	if err := s.t.ExecuteTemplate(rw, "about.gohtml", data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// Settings - display settings page
func (s *Service) Settings(rw http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(*user.AuthUser)

	data := struct {
		User *user.AuthUser
	}{
		User: usr,
	}

	log.Printf("data:%v+\n", data)

	if err := s.t.ExecuteTemplate(rw, "settings.gohtml", data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
