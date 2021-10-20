package handlers

import (
	"net/http"
	"photo-contest/business/data/contest"
	"photo-contest/business/data/photo"
	"photo-contest/business/data/user"
)

// ContestIndex - Contest Home
func (s *Service) ContestIndex(rw http.ResponseWriter, r *http.Request) {
}

// ContestRules - Contest Rules
func (s *Service) ContestRules(rw http.ResponseWriter, r *http.Request) {
}

// ContestPhotos - list contest photos
func (s *Service) ContestPhotos(rw http.ResponseWriter, r *http.Request) {
	contestID := 1
	usr := r.Context().Value("user").(*user.AuthUser)
	photoStore := photo.NewStore(s.log, s.db)
	contestStore := contest.NewStore(s.log, s.db)
	contestEntries, err := contestStore.QueryContestEntries(contestID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var thumbPhotos []PhotoInfo
	for _, contestEntry := range contestEntries {
		photo, err := photoStore.QueryByID(contestEntry.PhotoID)
		if err == nil {
			if !photo.Deleted {
				photoFile, _ := photoStore.QueryPhotoFile(photo.ID, "thumb")
				photoInfo := PhotoInfo{
					FilePath: photoFile.FilePath,
					PhotoId:  photo.ID,
					/*Title:          photo.Title,
					Description:      photo.Description,*/
					SubjectName:      contestEntry.SubjectName,
					SubjectAge:       contestEntry.SubjectAge,
					SubjectCountry:   contestEntry.SubjectCountry,
					SubjectOrigin:    contestEntry.SubjectOrigin,
					SubjectBiography: contestEntry.SubjectBiography,
					Location:         contestEntry.Location,
					ReleaseMimeType:  contestEntry.ReleaseMimeType,
					Status:           contestEntry.Status,
				}
				thumbPhotos = append(thumbPhotos, photoInfo)
			}
		}
	}
	formData := map[string]interface{}{
		"Photos": thumbPhotos,
		"User":   usr,
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Cache-Control", "no-cache")
	s.ExecuteTemplateWithBase(rw, formData, "gallery.gohtml")
}

// ContestAddPhoto - add a photo to a contest (req auth)
func (s *Service) ContestAddPhoto(rw http.ResponseWriter, r *http.Request) {
}
