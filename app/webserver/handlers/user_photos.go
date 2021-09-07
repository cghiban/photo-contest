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
		s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
	} else if r.Method == "POST" {
		formData["Message"] = "Unimplemented"
		s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
	}
}
