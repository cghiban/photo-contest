package handlers

import (
	"context"
	"fmt"
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

		s.log.Println("trying to find user w:", email, password)

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
			s.log.Println(err)
			formData["Message"] = err.Error()
			if err := s.t.ExecuteTemplate(rw, "register.gohtml", formData); err != nil {
				//http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			http.Error(rw, "Unable to sign user up", http.StatusInternalServerError)
			return
		} else {
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
		usr, err := userGroup.Authenticate(email, password)
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
				s.log.Printf("err = %+v\n", err)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			s.log.Println("all's good. will be redirecting to / now..")
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

// UserUpdateProfile - allows profile to be updated
func (s *Service) UserUpdateProfile(rw http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(*user.AuthUser)

	formData := map[string]interface{}{
		"name":           usr.Name,
		"email":          usr.Email,
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	if r.Method == "GET" {
		rw.Header().Add("Cache-Control", "no-cache")
		if err := s.t.ExecuteTemplate(rw, "profile.gohtml", formData); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		r.ParseForm()

		name := strings.Trim(r.Form.Get("name"), " ")
		email := strings.Trim(r.Form.Get("email"), " ")

		s.log.Println("trying to update user w:", email, name)

		uu := user.UpdateAuthUser{
			Name:  name,
			Email: email,
		}
		userGroup := user.NewStore(s.log, s.db)

		u, _ := userGroup.QueryByEmail(email)
		if u.ID != 0 && u.ID != usr.ID {
			formData["Message"] = "This email is already in use."
			if err := s.t.ExecuteTemplate(rw, "profile.gohtml", formData); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		_, err := userGroup.Update(usr.ID, uu)
		if err != nil {
			s.log.Println(err)
			formData["Message"] = err.Error()
			if err := s.t.ExecuteTemplate(rw, "profile.gohtml", formData); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(rw, r, "/profile", http.StatusFound)
	}
}

// UserUpdatePassword - allows profile to be updated
func (s *Service) UserUpdatePassword(rw http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(*user.AuthUser)

	formData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	if r.Method == "GET" {
		rw.Header().Add("Cache-Control", "no-cache")
		if err := s.t.ExecuteTemplate(rw, "password.gohtml", formData); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		r.ParseForm()

		old_password := strings.Trim(r.Form.Get("old_password"), " ")
		password := strings.Trim(r.Form.Get("password"), " ")
		password_confirm := strings.Trim(r.Form.Get("password_confirm"), " ")

		userGroup := user.NewStore(s.log, s.db)
		usr, err := userGroup.Authenticate(usr.Email, old_password)
		if err != nil {
			s.log.Println(err)
			formData["Message"] = "Incorrect former password"
			if err := s.t.ExecuteTemplate(rw, "password.gohtml", formData); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if usr != nil {
			s.log.Println("trying to update user password: ", usr.ID)

			up := user.UpdateAuthUserPass{
				Pass:        password,
				PassConfirm: password_confirm,
			}

			_, err := userGroup.UpdatePass(usr.ID, up)
			if err != nil {

				formData["Message"] = err.Error()
				if err := s.t.ExecuteTemplate(rw, "password.gohtml", formData); err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			http.Redirect(rw, r, "/", http.StatusFound)
		}
	}
}

// UserAuth provides middleware functions for authorizing users and setting the user
// in the request context.
type Auth struct {
	service *Service
}

// NewAuth constructs a Auth object store for checking logged in user.
func NewAuth(s *Service) Auth {
	return Auth{service: s}
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
