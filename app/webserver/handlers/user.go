package handlers

import (
	"context"
	"fmt"
	"net/http"
	"photo-contest/business/data/user"
	"photo-contest/foundation/utils"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/csrf"
)

// UserSignUp - handles user signup
func (s *Service) UserSignUp(rw http.ResponseWriter, r *http.Request) {
	states := utils.USStates()
	state_keys := utils.StateKeys(states)
	ethnicities := utils.Ethnicities()
	ethnicity_keys := utils.EthnicitiesKeys()
	genders := utils.Genders()
	gender_keys := utils.GenderKeys()
	sort.Sort(state_keys)
	formData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"states":         states,
		"state_keys":     state_keys,
		"ethnicities":    ethnicities,
		"ethnicity_keys": ethnicity_keys,
		"genders":        genders,
		"gender_keys":    gender_keys,
	}
	if r.Method == "GET" {
		s.ExecuteTemplateWithBase(rw, formData, "register.gohtml")
	} else if r.Method == "POST" {
		r.ParseForm()

		name := strings.TrimSpace(r.Form.Get("name"))
		email := strings.TrimSpace(r.Form.Get("email"))
		password := strings.TrimSpace(r.Form.Get("password"))
		password_confirm := strings.TrimSpace(r.Form.Get("password_confirm"))
		street := strings.TrimSpace(r.Form.Get("street"))
		city := strings.TrimSpace(r.Form.Get("city"))
		state := strings.TrimSpace(r.Form.Get("state"))
		zip := strings.TrimSpace(r.Form.Get("zip"))
		phone := strings.TrimSpace(r.Form.Get("phone"))
		age_string := strings.TrimSpace(r.Form.Get("age"))
		age, err := strconv.Atoi(age_string)
		if err != nil {
			formData["Message"] = "Age must be a number."
			s.ExecuteTemplateWithBase(rw, formData, "register.gohtml")
			return
		}
		if age < 0 || age > 255 {
			formData["Message"] = "Age must be a positive number."
			s.ExecuteTemplateWithBase(rw, formData, "register.gohtml")
			return
		}
		gender := strings.TrimSpace(r.Form.Get("gender"))
		ethnicities := r.Form["ethnicity"]
		ethnicity := ""
		for _, eth := range ethnicities {
			if utils.InStringSlice(ethnicity_keys, eth) {
				ethnicity = ethnicity + eth
				ethnicity = ethnicity + ";"
			}
		}
		ethnicity = strings.TrimSuffix(ethnicity, ";")
		other_ethnicity := strings.TrimSpace(r.Form.Get("other_ethnicity"))
		valid_state := false
		for _, s := range state_keys {
			if s == state {
				valid_state = true
			}
		}
		if !valid_state {
			formData["Message"] = "Invalid state."
			s.ExecuteTemplateWithBase(rw, formData, "register.gohtml")
			return
		}
		if !utils.InStringSlice(gender_keys, gender) {
			formData["Message"] = "Invalid gender."
			s.ExecuteTemplateWithBase(rw, formData, "register.gohtml")
			return
		}
		if password != password_confirm {
			formData["Message"] = "The passwords must match."
			s.ExecuteTemplateWithBase(rw, formData, "register.gohtml")
			return
		}

		s.log.Println("trying to find user w:", email, password)

		userGroup := user.NewStore(s.log, s.db)

		_, err = userGroup.QueryByEmail(email)
		if err == nil {
			formData["Message"] = "This email is already in use."
			s.ExecuteTemplateWithBase(rw, formData, "register.gohtml")
			return
		}

		newUser := user.NewAuthUser{
			Name:           name,
			Email:          email,
			Pass:           password,
			PassConfirm:    password_confirm,
			Street:         street,
			City:           city,
			State:          state,
			Zip:            zip,
			Phone:          phone,
			Age:            age,
			Gender:         gender,
			Ethnicity:      ethnicity,
			OtherEthnicity: other_ethnicity,
		}

		usr, err := userGroup.Create(newUser)
		fmt.Printf("usr = %+v\n", usr)
		fmt.Printf("err = %+v\t%T\n", err, err)
		if err != nil {
			s.log.Println(err)
			formData["Message"] = err.Error()
			s.ExecuteTemplateWithBaseNoServerError(rw, formData, "register.gohtml")
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
		s.ExecuteTemplateWithBase(rw, formData, "login.gohtml")
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
		s.ExecuteTemplateWithBase(rw, formData, "login.gohtml")
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
	states := utils.USStates()
	state_keys := utils.StateKeys(states)
	ethnicities := utils.Ethnicities()
	ethnicity_keys := utils.EthnicitiesKeys()
	genders := utils.Genders()
	gender_keys := utils.GenderKeys()
	sort.Sort(state_keys)
	formData := map[string]interface{}{
		"name":            usr.Name,
		"email":           usr.Email,
		"street":          usr.Street,
		"city":            usr.City,
		"state":           usr.State,
		"zip":             usr.Zip,
		"phone":           usr.Phone,
		"age":             usr.Age,
		"gender":          usr.Gender,
		"ethnicity":       usr.Ethnicity,
		"other_ethnicity": usr.OtherEthnicity,
		csrf.TemplateTag:  csrf.TemplateField(r),
		"User":            usr,
		"states":          states,
		"state_keys":      state_keys,
		"ethnicities":     ethnicities,
		"ethnicity_keys":  ethnicity_keys,
		"genders":         genders,
		"gender_keys":     gender_keys,
	}

	if r.Method == "GET" {
		rw.Header().Add("Cache-Control", "no-cache")
		s.ExecuteTemplateWithBase(rw, formData, "profile.gohtml")
	} else if r.Method == "POST" {
		r.ParseForm()

		name := strings.Trim(r.Form.Get("name"), " ")
		email := strings.Trim(r.Form.Get("email"), " ")
		street := strings.TrimSpace(r.Form.Get("street"))
		city := strings.TrimSpace(r.Form.Get("city"))
		state := strings.TrimSpace(r.Form.Get("state"))
		zip := strings.TrimSpace(r.Form.Get("zip"))
		phone := strings.TrimSpace(r.Form.Get("phone"))
		age_string := strings.TrimSpace(r.Form.Get("age"))
		age, err := strconv.Atoi(age_string)
		if err != nil {
			formData["Message"] = "Age must be a number."
			s.ExecuteTemplateWithBase(rw, formData, "profile.gohtml")
			return
		}
		if age < 0 || age > 255 {
			formData["Message"] = "Age must be a positive number."
			s.ExecuteTemplateWithBase(rw, formData, "profile.gohtml")
			return
		}
		gender := strings.TrimSpace(r.Form.Get("gender"))
		ethnicities := r.Form["ethnicity"]
		ethnicity := ""
		for _, eth := range ethnicities {
			if utils.InStringSlice(ethnicity_keys, eth) {
				ethnicity = ethnicity + eth
				ethnicity = ethnicity + ";"
			}
		}
		ethnicity = strings.TrimSuffix(ethnicity, ";")
		other_ethnicity := strings.TrimSpace(r.Form.Get("other_ethnicity"))
		valid_state := false
		for _, s := range state_keys {
			if s == state {
				valid_state = true
			}
		}
		if !valid_state {
			formData["Message"] = "Invalid state."
			s.ExecuteTemplateWithBase(rw, formData, "profile.gohtml")
			return
		}
		if !utils.InStringSlice(gender_keys, gender) {
			formData["Message"] = "Invalid gender."
			s.ExecuteTemplateWithBase(rw, formData, "profile.gohtml")
			return
		}
		s.log.Println("trying to update user w:", email, name)

		uu := user.UpdateAuthUser{
			Name:           name,
			Email:          email,
			Street:         street,
			City:           city,
			State:          state,
			Zip:            zip,
			Phone:          phone,
			Age:            age,
			Gender:         gender,
			Ethnicity:      ethnicity,
			OtherEthnicity: other_ethnicity,
		}
		userGroup := user.NewStore(s.log, s.db)

		u, _ := userGroup.QueryByEmail(email)
		if u.ID != 0 && u.ID != usr.ID {
			formData["Message"] = "This email is already in use."
			s.ExecuteTemplateWithBase(rw, formData, "profile.gohtml")
			return
		}

		_, err = userGroup.Update(usr.ID, uu)
		if err != nil {
			s.log.Println(err)
			formData["Message"] = err.Error()
			s.ExecuteTemplateWithBase(rw, formData, "profile.gohtml")
			return
		}

		http.Redirect(rw, r, "/profile", http.StatusFound)
	}
}

func (s *Service) UserForgotPassword(rw http.ResponseWriter, r *http.Request) {
	formData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	if r.Method == "GET" {
		rw.Header().Add("Cache-Control", "no-cache")
		s.ExecuteTemplateWithBase(rw, formData, "forgotpass.gohtml")
	} else if r.Method == "POST" {
		email := strings.TrimSpace(r.Form.Get("email"))
		userStore := user.NewStore(s.log, s.db)
		usr, err := userStore.QueryByEmail(email)
		if err != nil {
			s.ExecuteTemplateWithBase(rw, formData, "forgotpassredirect.gohtml")
			return
		}
		nr := user.NewResetPasswordEmail{
			UserID:    usr.ID,
			UpdatedBy: usr.Name,
		}
		pr, err := userStore.CreatePasswordReset(nr)
		if err != nil {
			formData["Message"] = "Unable to create password reset"
			s.ExecuteTemplateWithBase(rw, formData, "forgotpass.gohtml")
			return
		}
		_, _, err = utils.SendResetEmail(s.mailgunServer, usr.Email, pr.ResetID)
		if err != nil {
			fmt.Println(usr.Email, err)
			formData["message"] = "Unable to create password reset"
			s.ExecuteTemplateWithBase(rw, formData, "forgotpass.gohtml")
			return
		}
		s.ExecuteTemplateWithBase(rw, formData, "forgotpassredirect.gohtml")
	}
}

func (s *Service) UserResetPassword(rw http.ResponseWriter, r *http.Request) {
	formData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	userStore := user.NewStore(s.log, s.db)
	if r.Method == "GET" {
		codes, ok := r.URL.Query()["c"]
		code := ""
		if ok && len(codes[0]) > 0 {
			code = codes[0]
		}
		if code == "" {
			s.NotFoundHandler(rw, r)
			return
		}
		_, err := userStore.QueryPasswordResetByID(code)
		if err != nil {
			s.NotFoundHandler(rw, r)
			return
		}
		formData["Code"] = code
		rw.Header().Add("Cache-Control", "no-cache")
		s.ExecuteTemplateWithBase(rw, formData, "resetpass.gohtml")
	} else if r.Method == "POST" {
		code := strings.TrimSpace(r.Form.Get("code"))
		if code == "" {
			s.NotFoundHandler(rw, r)
			return
		}
		re, err := userStore.QueryPasswordResetByID(code)
		if err != nil {
			s.NotFoundHandler(rw, r)
			return
		}
		password := strings.TrimSpace(r.Form.Get("password"))
		password_confirm := strings.TrimSpace(r.Form.Get("password_confirm"))
		if password != password_confirm {
			formData["Message"] = "New passwords must match"
			s.ExecuteTemplateWithBase(rw, formData, "resetpass.gohtml")
			return
		}
		s.log.Println("trying to update user password: ", re.UserID)

		up := user.UpdateAuthUserPass{
			Pass:        password,
			PassConfirm: password_confirm,
		}

		_, err = userStore.UpdatePass(re.UserID, up)
		if err != nil {
			formData["Message"] = "Failed to update password"
			s.ExecuteTemplateWithBase(rw, formData, "resetpass.gohtml")
			return
		}
		er := user.ExpireResetPasswordEmail{
			ResetID: re.ResetID,
		}
		// If reset fails to be expired, it just will remain active until 24 hours have passed
		// Not a major deal, so no error checking needed
		userStore.ExpirePasswordReset(er)
		http.Redirect(rw, r, "/", http.StatusFound)
	}
}

// UserUpdatePassword - allows profile to be updated
func (s *Service) UserUpdatePassword(rw http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(*user.AuthUser)

	formData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"User":           usr,
	}

	if r.Method == "GET" {
		rw.Header().Add("Cache-Control", "no-cache")
		s.ExecuteTemplateWithBase(rw, formData, "password.gohtml")
	} else if r.Method == "POST" {
		r.ParseForm()

		old_password := strings.Trim(r.Form.Get("old_password"), " ")
		password := strings.Trim(r.Form.Get("password"), " ")
		password_confirm := strings.Trim(r.Form.Get("password_confirm"), " ")
		if password != password_confirm {
			formData["Message"] = "New passwords must match"
			s.ExecuteTemplateWithBase(rw, formData, "password.gohtml")
			return
		}
		userGroup := user.NewStore(s.log, s.db)
		usr, err := userGroup.Authenticate(usr.Email, old_password)
		if err != nil {
			s.log.Println(err)
			formData["Message"] = "Incorrect former password"
			s.ExecuteTemplateWithBase(rw, formData, "password.gohtml")
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
				formData["Message"] = "Failed to update password"
				s.ExecuteTemplateWithBase(rw, formData, "password.gohtml")
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

func (a *Auth) RequireAdmin(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmp := r.Context().Value("user")
		if tmp == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		usr, ok := tmp.(*user.AuthUser)
		if !ok {
			// Whatever was set in the user key isn't a user, so we probably need to
			// sign in.
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if usr.PermissionLevel < 2 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
}
