package handlers

import (
	"log"
	"net/http"
	"photo-contest/business/data"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
)

// Service data struct
type Service struct {
	log   *log.Logger
	store *data.DataStore
	//readers *string
	//session *sqlitestore.SqliteStore
	session *sessions.CookieStore
	t       *template.Template
}

// NewService initializes a new Serivice
func NewService(l *log.Logger, store *data.DataStore, sessionKey string) *Service {
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

	//sessStore.Options = &sessions.Options{HttpOnly: true}

	sessStore.Options = &sessions.Options{
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7 * 86400,
	}

	return &Service{log: l, store: store, t: templates, session: sessStore}
}

// About - about this site
func (s *Service) About(rw http.ResponseWriter, r *http.Request) {
	var user *data.AuthUser
	userV := r.Context().Value("user")
	if userV != nil {
		user = userV.(*data.AuthUser)
	}
	data := struct {
		User *data.AuthUser
	}{
		User: user,
	}
	if err := s.t.ExecuteTemplate(rw, "about.gohtml", data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// Settings - display settings page
func (s *Service) Settings(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*data.AuthUser)

	data := struct {
		User *data.AuthUser
	}{
		User: user,
	}

	log.Printf("data:%v+\n", data)

	if err := s.t.ExecuteTemplate(rw, "settings.gohtml", data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
