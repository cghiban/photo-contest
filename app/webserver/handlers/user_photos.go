package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"io"
	"log"
	"os"
	"path"

	"photo-contest/business/data/contest"
	"photo-contest/business/data/photo"
	"photo-contest/business/data/user"

	"github.com/gorilla/csrf"
	"github.com/h2non/bimg"
)

type PhotoInfo struct {
	FilePath string
	/*Title          string
	Description      string*/
	Status           string
	PhotoId          string
	SubjectName      string
	SubjectAge       string
	SubjectCountry   string
	SubjectOrigin    string
	SubjectBiography string
	Location         string
	ReleaseMimeType  string
	EntryId          int
}

// UserPhotos - lists user photos (req auth)
func (s *Service) UserPhotos(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		usr := r.Context().Value("user").(*user.AuthUser)
		photoStore := photo.NewStore(s.log, s.db)
		contestStore := contest.NewStore(s.log, s.db)
		userPhotos, err := photoStore.QueryByOwnerID(usr.ID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		var thumbPhotos []PhotoInfo
		for _, photo := range userPhotos {
			if !photo.Deleted {
				photoFile, _ := photoStore.QueryPhotoFile(photo.ID, "thumb")
				contestEntry, _ := contestStore.QueryContestEntryByPhotoId(photo.ID)
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
		formData := map[string]interface{}{
			"Photos": thumbPhotos,
			"User":   usr,
			"CSRF":   csrf.Token(r),
		}
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Cache-Control", "no-cache")
		s.ExecuteTemplateWithBase(rw, formData, "submissions.gohtml")
	}
}

func handleFileUpload(r *http.Request, photoID string) error {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(11 << 20)
	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, fileInfo, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", fileInfo.Filename)
	fmt.Printf("File Size: %+v\n", fileInfo.Size)
	fmt.Printf("Content Type: %+v\n", fileInfo.Header["Content-Type"][0])
	if fileInfo.Header["Content-Type"][0] != "image/jpeg" && fileInfo.Header["Content-Type"][0] != "image/png" {
		return fmt.Errorf("Invalid Content Type: %s", fileInfo.Header["Content-Type"][0])
	}

	// copy all of the contents of our uploaded file into a
	// new file
	new_file, err := os.Create(fmt.Sprintf("tmp/photo-%s-original.jpg", photoID))
	if err != nil {
		return err
	}
	_, err = io.Copy(new_file, file)
	return err
}

func imageTooBig(photoID string) bool {
	// Get the full path to the full size image and the full path to the location to save the thumbnail
	file := fmt.Sprintf("tmp/photo-%s-original.jpg", photoID)
	// Read the image from its location
	buffer, err := bimg.Read(file)
	if err != nil {
		fmt.Println("Image too big check; Buffer: " + err.Error())
		return false
	}
	// Check image size
	size, err := bimg.NewImage(buffer).Size()
	if err != nil {
		fmt.Println("Image too big check; Size: " + err.Error())
		return false
	}
	// Make sure the image is not too large
	if size.Width > 2000 || size.Height > 2000 {
		return true
	}
	// Write the thumbnail to the proper location
	return false
}

func handleModelReleaseUpload(r *http.Request, photoID string) (string, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `release`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, fileInfo, err := r.FormFile("release")
	if err != nil {
		return "", err
	}
	defer file.Close()
	modelReleaseTypes := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
	}
	mimeType := fileInfo.Header["Content-Type"][0]
	if !modelReleaseTypes[mimeType] {
		return "", fmt.Errorf("Invalid Content Type: %s", mimeType)
	}

	// copy all of the contents of our uploaded file into a
	// new file
	extension := "pdf"
	if mimeType == "image/png" {
		extension = "png"
	} else if mimeType == "image/jpeg" {
		extension = "jpg"
	}
	new_file, err := os.Create(fmt.Sprintf("tmp/release-%s.%s", photoID, extension))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(new_file, file)
	return mimeType, err
}

// Create a resized image from the full size image
func makeThumbnail(id string, s string, pixels uint16, l *log.Logger) error {
	// Get the full path to the full size image and the full path to the location to save the thumbnail
	fullPathToThumbnail := fmt.Sprintf("tmp/photo-%s-%s.jpg", id, s)
	file := fmt.Sprintf("tmp/photo-%s-original.jpg", id)
	// The full path to the directory that needs to be created to save the thumbnail properly
	fullDirectory := path.Dir(fullPathToThumbnail)
	// Make all layers of directory that are needed, assuming they do not exist; when they do exist, this function does nothing
	err := os.MkdirAll(fullDirectory, os.ModePerm)
	if err != nil {
		return err
	}
	// Read the image from its location
	buffer, err := bimg.Read(file)
	if err != nil {
		return err
	}
	// Use the bimg Thumbnail function, which resizes and crops the image, with a quality of 95%, using libvips
	newImage, err := bimg.NewImage(buffer).Thumbnail(int(pixels))
	if err != nil {
		return err
	}
	size, err := bimg.NewImage(newImage).Size()
	if err != nil {
		return err
	}
	// Make sure the crop function created an image of the proper size; enlarge it if not
	if size.Width != int(pixels) || size.Height != int(pixels) {
		newImage, err = bimg.NewImage(buffer).EnlargeAndCrop(int(pixels), int(pixels))
		if err != nil {
			return err
		}
	}
	// Write the thumbnail to the proper location
	return bimg.Write(fullPathToThumbnail, newImage)
}

// UserPhotoUpload - upload user photos (req auth)
func (s *Service) UserPhotoUpload(rw http.ResponseWriter, r *http.Request) {
	// on GET display the form
	// on POST handle file upload
	photoStore := photo.NewStore(s.log, s.db)
	contestStore := contest.NewStore(s.log, s.db)
	usr := r.Context().Value("user").(*user.AuthUser)
	ids, ok := r.URL.Query()["id"]
	id := ""
	if ok && len(ids[0]) > 0 {
		id = ids[0]
	}
	contestID := 1
	formData := map[string]interface{}{
		csrf.TemplateTag:   csrf.TemplateField(r),
		"User":             usr,
		"SubjectName":      "",
		"SubjectAge":       "",
		"SubjectCountry":   "",
		"SubjectOrigin":    "",
		"Location":         "",
		"SubjectBiography": "",
		"Signature":        "",
		"RealName":         usr.Name,
		"Edit":             0,
		"ID":               id,
	}
	var photoTemplate photo.Photo
	if id != "" {
		photo, err := photoStore.QueryByID(id)
		if err == nil {
			if photo.OwnerID == usr.ID {
				photoTemplate = photo
			}
		}
	}
	var contestEntry contest.ContestEntry
	if photoTemplate.ID != "" {
		contestEntry, _ = contestStore.QueryContestEntryByPhotoId(photoTemplate.ID)
	}
	if contestEntry.Status == "active" {
		formData["SubjectName"] = contestEntry.SubjectName
		formData["SubjectAge"] = contestEntry.SubjectAge
		formData["SubjectCountry"] = contestEntry.SubjectCountry
		formData["SubjectOrigin"] = contestEntry.SubjectOrigin
		formData["Location"] = contestEntry.Location
		formData["SubjectBiography"] = contestEntry.SubjectBiography
		formData["Signature"] = usr.Name
		formData["Edit"] = 1
	}
	if r.Method == "GET" {
		rw.Header().Add("Cache-Control", "no-cache")
		s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
	} else if r.Method == "POST" {
		usr := r.Context().Value("user").(*user.AuthUser)
		title := ""
		description := ""
		subject_name := strings.TrimSpace(r.Form.Get("sname"))
		subject_age := strings.TrimSpace(r.Form.Get("sage"))
		subject_country := strings.TrimSpace(r.Form.Get("scountry"))
		subject_origin := strings.TrimSpace(r.Form.Get("sorigin"))
		location := strings.TrimSpace(r.Form.Get("location"))
		subject_biography := strings.TrimSpace(r.Form.Get("sbiography"))
		signature := strings.TrimSpace(r.Form.Get("signature"))
		formData["SubjectName"] = subject_name
		formData["SubjectAge"] = subject_age
		formData["SubjectCountry"] = subject_country
		formData["SubjectOrigin"] = subject_origin
		formData["Location"] = location
		formData["SubjectBiography"] = subject_biography
		formData["Signature"] = signature
		r.ParseMultipartForm(10 << 20)
		//title := strings.TrimSpace(r.Form.Get("title"))
		//description := strings.TrimSpace(r.Form.Get("description"))
		if usr.Name != signature {
			formData["Message"] = "Digital Signature (" + signature + ") must exactly match the name in your account (" + usr.Name + "). If the name in your account is not a legal name, update your profile."
			s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
			return
		}
		if len(subject_biography) < 250 || len(subject_biography) > 500 {
			formData["Message"] = "Biography must be between 250 and 500 characters."
			s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
			return
		}
		if formData["Edit"] == 1 {
			fuce := contest.FullyUpdateContestEntry{
				EntryID:          contestEntry.EntryID,
				UpdatedBy:        usr.Name,
				SubjectName:      subject_name,
				SubjectAge:       subject_age,
				SubjectCountry:   subject_country,
				SubjectOrigin:    subject_origin,
				Location:         location,
				SubjectBiography: subject_biography,
			}
			_, err := contestStore.FullyUpdateEntry(fuce)
			if err != nil {
				formData["Message"] = "Unable to update entry."
				s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
				return
			}
			formData["Message"] = "Updated the contest entry"
		} else {
			userPhotos, err := photoStore.QueryByOwnerID(usr.ID)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			userSubmittedPhotos := 0
			for _, photo := range userPhotos {
				ce, err := contestStore.QueryContestEntryByPhotoId(photo.ID)
				if err != nil {
					userSubmittedPhotos = userSubmittedPhotos + 1
					continue
				}
				if ce.Status != "withdrawn" {
					userSubmittedPhotos = userSubmittedPhotos + 1
				}
			}
			if userSubmittedPhotos >= 3 {
				formData["Message"] = "There is a maximum of three submissions per user and you've already submitted three."
				s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
				return
			}
			np := photo.NewPhoto{
				OwnerID:     usr.ID,
				Title:       title,
				Description: description,
				UpdatedBy:   usr.Name,
			}
			pht, err := photoStore.Create(np)
			if err != nil {
				formData["Message"] = "Could not upload photo"
				s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
				return
			}
			mimeType, err := handleModelReleaseUpload(r, pht.ID)
			if err != nil {
				formData["Message"] = "Could not upload model release form"
				s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
				return
			}
			nce := contest.NewContestEntry{
				ContestID:        contestID,
				PhotoID:          pht.ID,
				Status:           "active",
				UpdatedBy:        usr.Name,
				SubjectName:      subject_name,
				SubjectAge:       subject_age,
				SubjectCountry:   subject_country,
				SubjectOrigin:    subject_origin,
				Location:         location,
				ReleaseMimeType:  mimeType,
				SubjectBiography: subject_biography,
			}
			_, err = contestStore.CreateContestEntry(nce)
			if err != nil {
				err = photoStore.Delete(pht.ID)
				if err != nil {
					fmt.Println("Could not delete: " + err.Error())
				}
				formData["Message"] = "Could not create contest entry"
				s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
				return
			}
			npf := photo.NewPhotoFile{
				PhotoID:   pht.ID,
				FilePath:  fmt.Sprintf("/tmp/photo-%s-original.jpg", pht.ID),
				Size:      "original",
				UpdatedBy: usr.Name,
			}
			_, err = photoStore.CreateFile(npf)
			if err != nil {
				err = photoStore.Delete(pht.ID)
				if err != nil {
					fmt.Println("Could not delete: " + err.Error())
				}
				formData["Message"] = "Could not upload photo"
				s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
				return
			}
			err = handleFileUpload(r, pht.ID)
			if err != nil {
				err = photoStore.Delete(pht.ID)
				if err != nil {
					fmt.Println("Could not delete: " + err.Error())
				} else {
					fmt.Println("Deleted?")
				}
				formData["Message"] = "Could not upload photo"
				s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
				return
			}
			if imageTooBig(pht.ID) {
				err = photoStore.Delete(pht.ID)
				if err != nil {
					fmt.Println("Could not delete: " + err.Error())
				}
				formData["Message"] = "Image must be no more than 2000 pixels along any dimension"
				s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
				return
			}
			photoSizes := []string{"thumb", "small", "medium", "large"}
			for _, size := range photoSizes {
				npf := photo.NewPhotoFile{
					PhotoID:   pht.ID,
					FilePath:  fmt.Sprintf("/tmp/photo-%s-%s.jpg", pht.ID, size),
					Size:      size,
					UpdatedBy: usr.Name,
				}
				phf, err := photoStore.CreateFile(npf)
				if err != nil {
					err = photoStore.Delete(pht.ID)
					if err != nil {
						fmt.Println("Could not delete: " + err.Error())
					}
					formData["Message"] = "Could not upload photo"
					s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
					return
				}
				err = makeThumbnail(pht.ID, size, phf.Width, s.log)
				if err != nil {
					err = photoStore.Delete(pht.ID)
					if err != nil {
						fmt.Println("Could not delete: " + err.Error())
					}
					formData["Message"] = "Could not upload photo"
					s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
					return
				}
			}
			formData["Message"] = "Uploaded the image"
		}
		s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
	}
}

func (s *Service) UserWithdrawPhoto(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusForbidden)
		rw.Header().Set("Content-Type", "application/json")
		return
	}
	usr := r.Context().Value("user").(*user.AuthUser)
	jsonEnc := json.NewEncoder(rw)
	response := map[string]string{
		"status":  "error",
		"message": "",
	}
	photoStore := photo.NewStore(s.log, s.db)
	contestStore := contest.NewStore(s.log, s.db)
	photoId := r.Form.Get("photoId")
	photo, err := photoStore.QueryByID(photoId)
	if err != nil {
		response["message"] = "No photo found"
		jsonEnc.Encode(response)
		return
	}
	if photo.OwnerID != usr.ID {
		response["message"] = "No photo found"
		jsonEnc.Encode(response)
		return
	}
	entry, err := contestStore.QueryContestEntryByPhotoId(photoId)
	if err != nil {
		response["message"] = "No contest entry found"
		jsonEnc.Encode(response)
		return
	}
	uce := contest.UpdateContestEntry{
		EntryID: entry.EntryID,
		Status:  "withdrawn",
	}
	_, err = contestStore.UpdateEntry(uce)
	if err != nil {
		response["message"] = "Can't withdraw contest entry"
		jsonEnc.Encode(response)
		return
	}
	response["status"] = "success"
	jsonEnc.Encode(response)
}
