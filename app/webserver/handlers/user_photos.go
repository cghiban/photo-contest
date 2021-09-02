package handlers

import "net/http"

// UserPhotos - lists user photos (req auth)
func (s *Service) UserPhotos(rw http.ResponseWriter, r *http.Request) {
}

// UserPhotoUpload - upload user photos (req auth)
func (s *Service) UserPhotoUpload(rw http.ResponseWriter, r *http.Request) {
	// on GET display the form
	// on POST handle file upload
}
