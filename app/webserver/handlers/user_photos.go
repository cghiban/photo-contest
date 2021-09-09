package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/json"
	"io"
	"os"
	"path"

	"photo-contest/business/data/photo"
	"photo-contest/business/data/user"

	"github.com/gorilla/csrf"
)

type PhotoInfo struct {
	FilePath    string
	Title       string
	Description string
	PhotoId     string
}

type ThumbnailResponse struct {
	Status  string
	Message string
}

// UserPhotos - lists user photos (req auth)
func (s *Service) UserPhotos(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		usr := r.Context().Value("user").(*user.AuthUser)
		photoStore := photo.NewStore(s.log, s.db)
		userPhotos, err := photoStore.QueryByOwnerID(usr.ID)
		var thumbPhotos []PhotoInfo
		for _, photo := range userPhotos {
			if !photo.Deleted {
				photoFile, _ := photoStore.QueryPhotoFile(photo.ID, "thumb")
				photoInfo := PhotoInfo{FilePath: photoFile.FilePath, PhotoId: photo.ID, Title: photo.Title, Description: photo.Description}
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
	err = os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		return err
	}
	new_file, err := os.Create(fmt.Sprintf("tmp/photo-%s-original.jpg", photoID))
	if err != nil {
		return err
	}
	_, err = io.Copy(new_file, file)
	return err
}

// Create a resized image from the full size image
func makeThumbnail(hostName string, id string, s string, pixels uint16) error {
	// Get the full path to the full size image and the full path to the location to save the thumbnail
	fullPathToThumbnail := fmt.Sprintf("tmp/photo-%s-%s.jpg", id, s)
	fileURL := fmt.Sprintf("%s/tmp/photo-%s-original.jpg", hostName, id)
	// The full path to the directory that needs to be created to save the thumbnail properly
	fullDirectory := path.Dir(fullPathToThumbnail)
	// Make all layers of directory that are needed, assuming they do not exist; when they do exist, this function does nothing
	err := os.MkdirAll(fullDirectory, os.ModePerm)
	if err != nil {
		return err
	}
	// Read the image from its location
	requestParams := make(map[string]interface{})
	requestParams["image"] = fileURL
	requestParams["pixels"] = pixels
	jsonValue, err := json.Marshal(requestParams)
	if err != nil {
		return err
	}
	resp, err := http.Post(hostName[0:strings.LastIndex(hostName, ":")]+":8000", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	var tr ThumbnailResponse
	var newImage []byte
	if resp.Header.Get("Content-Type") == "application/json" {
		json.NewDecoder(resp.Body).Decode(&tr)
		if tr.Status == "error" {
			return fmt.Errorf("%s", tr.Message)
		}
	} else {
		newImage, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	// Write the thumbnail to the proper location
	return os.WriteFile(fullPathToThumbnail, newImage, os.ModePerm)
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
		r.ParseMultipartForm(10 << 20)
		title := r.Form.Get("title")
		description := r.Form.Get("description")
		photoStore := photo.NewStore(s.log, s.db)
		np := photo.NewPhoto{
			OwnerID:     usr.ID,
			Title:       title,
			Description: description,
			UpdatedBy:   usr.Name,
		}
		pht, err := photoStore.Create(np)
		if err != nil {
			formData["Message"] = "Could not upload photo"
			s.log.Println(err.Error())
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
			s.log.Println(err.Error())
			err = photoStore.Delete(pht.ID)
			if err != nil {
				s.log.Println(err.Error())
				fmt.Println("Could not delete: " + err.Error())
			}
			formData["Message"] = "Could not upload photo"
			s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
			return
		}
		err = handleFileUpload(r, pht.ID)
		if err != nil {
			s.log.Println(err.Error())
			err = photoStore.Delete(pht.ID)
			if err != nil {
				s.log.Println(err.Error())
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
				s.log.Println(err.Error())
				err = photoStore.Delete(pht.ID)
				if err != nil {
					s.log.Println(err.Error())
					fmt.Println("Could not delete: " + err.Error())
				}
				formData["Message"] = "Could not upload photo"
				s.ExecuteTemplateWithBase(rw, formData, "photo.gohtml")
				return
			}
			hostName := ""
			if r.TLS == nil {
				hostName = "http://" + r.Host
			} else {
				hostName = "https://" + r.Host
			}
			err = makeThumbnail(hostName, pht.ID, size, phf.Width)
			if err != nil {
				s.log.Println(err.Error())
				err = photoStore.Delete(pht.ID)
				if err != nil {
					s.log.Println(err.Error())
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
