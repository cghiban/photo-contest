package handlers

import (
	"encoding/json"
	"net/http"
	"photo-contest/business/data/contest"
	"photo-contest/business/data/photo"
	"photo-contest/business/data/user"
	"strconv"

	"github.com/gorilla/csrf"
)

// ContestIndex - Contest Home
func (s *Service) ContestIndex(rw http.ResponseWriter, r *http.Request) {
}

// ContestRules - Contest Rules
func (s *Service) ContestRules(rw http.ResponseWriter, r *http.Request) {
}

// ContestPhotos - list contest photos
func (s *Service) ContestPhotos(rw http.ResponseWriter, r *http.Request) {
	statuses, ok := r.URL.Query()["status"]
	statusFilter := ""
	if ok && len(statuses[0]) > 0 {
		statusFilter = statuses[0]
	}
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
		if contestEntry.Status != "withdrawn" {
			photo, err := photoStore.QueryByID(contestEntry.PhotoID)
			if err == nil {
				if !photo.Deleted {
					if statusFilter == "" || contestEntry.Status == statusFilter {
						photoFile, _ := photoStore.QueryPhotoFile(photo.ID, "thumb")
						photoInfo := PhotoInfo{
							FilePath: photoFile.FilePath,
							PhotoId:  photo.ID,
							/*Title:          photo.Title,
							Descristringsption:      photo.Description,*/
							EntryId:          contestEntry.EntryID,
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
		}
	}
	formData := map[string]interface{}{
		"CSRF":        csrf.Token(r),
		"Photos":      thumbPhotos,
		"User":        usr,
		"FullContest": true,
		"Status":      statusFilter,
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Cache-Control", "no-cache")
	s.ExecuteTemplateWithBase(rw, formData, "gallery.gohtml")
}

func (s *Service) ContestEntryUpdateStatus(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusForbidden)
		rw.Header().Set("Content-Type", "application/json")
		return
	}
	jsonEnc := json.NewEncoder(rw)
	response := map[string]string{
		"status":  "error",
		"message": "",
	}
	entryIdStr := r.Form.Get("entryId")
	entryId, err := strconv.Atoi(entryIdStr)
	if err != nil {
		response["message"] = err.Error()
		jsonEnc.Encode(response)
		return
	}
	status := r.Form.Get("status")
	contestStore := contest.NewStore(s.log, s.db)
	uce := contest.UpdateContestEntry{
		EntryID: entryId,
		Status:  status,
	}
	_, err = contestStore.UpdateEntry(uce)
	if err != nil {
		response["message"] = err.Error()
	}
	response["status"] = "success"
	jsonEnc.Encode(response)
}

// ContestAddPhoto - add a photo to a contest (req auth)
func (s *Service) ContestAddPhoto(rw http.ResponseWriter, r *http.Request) {
}
