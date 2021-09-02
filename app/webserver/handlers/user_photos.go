package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// UserPhotos - lists user photos (req auth)
func (s *Service) UserPhotos(rw http.ResponseWriter, r *http.Request) {
}

// UserPhotoUpload - upload user photos (req auth)
func (s *Service) UserPhotoUpload(rw http.ResponseWriter, r *http.Request) {
	// on GET display the form
	// on POST handle file upload
	formData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"User":           "Exists",
	}

	if r.Method == "GET" {
		rw.Header().Add("Cache-Control", "no-cache")
		if _, err := s.t.ParseFiles("var/templates/base.gohtml", "var/templates/photo.gohtml"); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		if err := s.t.ExecuteTemplate(rw, "base", formData); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		formData["Message"] = "Unimplemented"
		if _, err := s.t.ParseFiles("var/templates/base.gohtml", "var/templates/photo.gohtml"); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		if err := s.t.ExecuteTemplate(rw, "base", formData); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}
