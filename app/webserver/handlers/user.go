package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"photo-contest/business/data/user"
	"photo-contest/foundation/database"
	"strings"

	"github.com/gorilla/csrf"
)

// UserSignUp - handles user signup
func (s *Service) UserSignUp(rw http.ResponseWriter, r *http.Request) {

	formData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	if r.Method == "GET" {
		if err := s.t.ExecuteTemplate(rw, "register.gohtml", formData); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		r.ParseForm()

		name := strings.Trim(r.Form.Get("name"), " ")
		email := strings.Trim(r.Form.Get("email"), " ")
		password := strings.Trim(r.Form.Get("password"), " ")
		password_confirm := strings.Trim(r.Form.Get("password_confirm"), " ")

		log.Println("trying to find user w:", email, password)

		userGroup := user.NewStore(s.log, s.db)

		_, err := userGroup.QueryByEmail(email)
		s.log.Println("from GetUser:", err)
		if err != nil && err != database.ErrNotFound {
			formData["Message"] = "This email is already in use."
			if err := s.t.ExecuteTemplate(rw, "register.gohtml", formData); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		newUser := user.NewAuthUser{
			Name:        name,
			Email:       email,
			Pass:        password,
			PassConfirm: password_confirm,
		}

		usr, err := userGroup.Create(newUser)
		fmt.Printf("usr = %+v\n", usr)
		fmt.Printf("err = %+v\t%T\n", err, err)
		if err != nil {
			log.Println(err)
			formData["Message"] = err.Error()
			if err := s.t.ExecuteTemplate(rw, "register.gohtml", formData); err != nil {
				//http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			//http.Error(rw, "Unable to sign user up", http.StatusInternalServerError)
			return
		} else {
			//s.log.Printf("user: %#v", usr)
			http.Redirect(rw, r, "/login", http.StatusFound)
		}
	}
}

func (s *Service) UserLogIn(rw http.ResponseWriter, r *http.Request) {

	formData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	if r.Method == "GET" {

		rw.Header().Add("Cache-Control", "no-cache")
		if err := s.t.ExecuteTemplate(rw, "login.gohtml", formData); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {

		if err := r.ParseForm(); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		email := strings.Trim(r.Form.Get("email"), " ")
		password := strings.Trim(r.Form.Get("password"), " ")

		userGroup := user.NewStore(s.log, s.db)
		//if err != nil {
		//	//http.Error(rw, err.Error(), http.StatusInternalServerError)
		//} else
		usr, err := userGroup.Authenticate(email, password)
		//log.Printf("usr = %+v\n", usr)
		//log.Printf("err = %+v\n", err)
		if err == nil && usr != nil {
			session, err := s.session.Get(r, "session")
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			session.Values["logged_in"] = true
			session.Values["user_id"] = usr.ID
			session.Values["name"] = usr.Name

			err = session.Save(r, rw)
			if err != nil {
				log.Printf("err = %+v\n", err)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Println("all's good. will be redirecting to / now..")
			http.Redirect(rw, r, "/", http.StatusFound)
			return
		}

		formData["Message"] = "Invalid email or password!"
		if err := s.t.ExecuteTemplate(rw, "login.gohtml", formData); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

// UserLogOut - clears the session
func (s *Service) UserLogOut(rw http.ResponseWriter, r *http.Request) {

	session, err := s.session.Get(r, "session")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["logged_in"] = false
	session.Options.MaxAge = -1

	err = session.Save(r, rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusFound)
}

// UserAuth provides middleware functions for authorizing users and setting the user
// in the request context.
type Auth struct {
	service *Service
}

// NewAuth constructs a Auth object store for checking logged in user.
func NewAuth(s *Service) Auth {
	return Auth{
		service: s,
	}
}

// UserViaSession will retrieve the current user set by the session cookie
// and set it in the request context. UserViaSession will NOT redirect
// to the sign in page if the user is not found. That is left for the
// RequireUser method to handle so that some pages can optionally have
// access to the current user.
func (a *Auth) UserViaSession(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		session, err := a.service.session.Get(r, "session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		//a.Service.log.Printf("logged_in: %v\t%T", session.Values["logged_in"], session.Values["logged_in"])
		if session.Values["logged_in"] != true {
			next.ServeHTTP(w, r)
			return
		}

		user_id, _ := session.Values["user_id"].(int)

		userGroup := user.NewStore(a.service.log, a.service.db)
		usr, err := userGroup.QueryByID(user_id)
		if err != nil {
			// If you want you can retain the original functionality to call
			// http.Error if any error aside from app.ErrNotFound is returned,
			// but I find that most of the time we can continue on and let later
			// code error if it requires a user, otherwise it can continue without
			// the user.
			next.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "user", &usr))
		next.ServeHTTP(w, r)
	}
}

// RequireUser will verify that a user is set in the request context. It if is
// set correctly, the next handler will be called, otherwise it will redirect
// the user to the sign in page and the next handler will not be called.
func (a *Auth) RequireUser(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmp := r.Context().Value("user")
		if tmp == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if _, ok := tmp.(*user.AuthUser); !ok {
			// Whatever was set in the user key isn't a user, so we probably need to
			// sign in.
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
}
