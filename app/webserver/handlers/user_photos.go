package handlers

import (
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
			"User":   "Exists",
		}
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Cache-Control", "no-cache")
		s.ExecuteTemplateWithBase(rw, formData, "gallery.gohtml")
	}
}

func handleFileUpload(r *http.Request, photoID string) error {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
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
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}
	mimeType := fileInfo.Header["Content-Type"][0]
	if !modelReleaseTypes[mimeType] {
		return "", fmt.Errorf("Invalid Content Type: %s", mimeType)
	}

	// copy all of the contents of our uploaded file into a
	// new file
	extension := "docx"
	if mimeType == "application/pdf" {
		extension = "pdf"
	} else if mimeType == "application/msword" {
		extension = "doc"
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
	formData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"User":           "Exists",
	}

	if r.Method == "GET" {
		rw.Header().Add("Cache-Control", "no-cache")
		s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
	} else if r.Method == "POST" {
		usr := r.Context().Value("user").(*user.AuthUser)
		photoStore := photo.NewStore(s.log, s.db)
		userPhotos, err := photoStore.QueryByOwnerID(usr.ID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(userPhotos) >= 3 {
			formData["Message"] = "There is a maximum of three submissions per user and you've already submitted three."
			s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
			return
		}
		r.ParseMultipartForm(10 << 20)
		//title := strings.TrimSpace(r.Form.Get("title"))
		//description := strings.TrimSpace(r.Form.Get("description"))
		title := ""
		description := ""
		subject_name := strings.TrimSpace(r.Form.Get("sname"))
		subject_age := strings.TrimSpace(r.Form.Get("sage"))
		subject_country := strings.TrimSpace(r.Form.Get("scountry"))
		subject_origin := strings.TrimSpace(r.Form.Get("sorigin"))
		location := strings.TrimSpace(r.Form.Get("location"))
		subject_biography := strings.TrimSpace(r.Form.Get("sbiography"))
		signature := strings.TrimSpace(r.Form.Get("signature"))
		if usr.Name != signature {
			formData["Message"] = "Digital Signature (" + signature + ") must exactly match the name in your account (" + usr.Name + "). If the name in your account is not a legal name, update your profile."
			s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
			return
		}
		contestStore := contest.NewStore(s.log, s.db)
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
			ContestID:        1,
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
		s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
	}
}
