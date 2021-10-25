package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"photo-contest/business/data/user"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/mailgun/mailgun-go/v4"
)

type Configuration struct {
	MailgunAPIKey string `json:"mailgun_api_key"`
	MailgunDomain string `json:"mailgun_domain"`
}

// Service data struct
type Service struct {
	C             Configuration
	log           *log.Logger
	db            *sqlx.DB
	session       *sessions.CookieStore
	t             *template.Template
	mailgunServer *mailgun.MailgunImpl
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
		"replaceEnd": func(input, from, to string) string { return strings.TrimSuffix(input, from) + to },
		"dateISOish": func(t time.Time) string { return t.Format("2006-01-02 3:04pm") },
		"substringInStringWithSeparator": func(full_string, substring, separator string) bool {
			substrs := strings.Split(full_string, separator)
			for _, s := range substrs {
				if s == substring {
					return true
				}
			}
			return false
		},
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

	f, err := os.Open("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	config := Configuration{}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}
	mg := mailgun.NewMailgun(config.MailgunDomain, config.MailgunAPIKey)
	return &Service{log: l, db: db, t: templates, session: sessStore, mailgunServer: mg, C: config}
}

func (s *Service) ExecuteTemplateWithBase(rw http.ResponseWriter, data interface{}, fileName string) {
	if err := s.t.ExecuteTemplate(rw, fileName, data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Service) ExecuteTemplateWithBaseNoServerError(rw http.ResponseWriter, data interface{}, fileName string) {
	if err := s.t.ExecuteTemplate(rw, fileName, data); err != nil {
		//http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Service) NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
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
	rw.WriteHeader(http.StatusNotFound)
	s.ExecuteTemplateWithBase(rw, data, "404.gohtml")
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
	s.ExecuteTemplateWithBase(rw, data, "index.gohtml")
}

// About - about this site
func (s *Service) About(rw http.ResponseWriter, r *http.Request) {
	var usr *user.AuthUser
	userV := r.Context().Value("user")
	if userV != nil {
		usr = userV.(*user.AuthUser)
	}
	data := struct {
		User    *user.AuthUser
		Message string
	}{
		User: usr,
	}
	s.ExecuteTemplateWithBase(rw, data, "about.gohtml")
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

	s.ExecuteTemplateWithBase(rw, data, "settings.gohtml")
}
