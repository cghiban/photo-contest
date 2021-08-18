package photo

import (
	"fmt"
	"log"
	"photo-contest/business/data/tests"
	"photo-contest/business/sys/validate"
	"photo-contest/foundation/database"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Store manages the set of API's for photo access.
type Store struct {
	log *log.Logger
	db  *sqlx.DB
}

// NewStore constructs a photo store for api access.
func NewStore(log *log.Logger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// Create - add new photo into db
func (s Store) Create(np NewPhoto) (Photo, error) {

	if err := validate.Check(np); err != nil {
		return Photo{}, errors.Wrap(err, "validating data")
	}

	now := time.Now()

	pht := Photo{
		ID:          validate.GenerateID(),
		OwnerID:     np.OwnerID,
		Title:       np.Title,
		Description: np.Description,
		CreatedOn:   now,
		UpdatedOn:   now,
		UpdatedBy:   np.UpdatedBy,
	}

	const query = `
	INSERT INTO photos
		(photo_id, owner_id, title, description, deleted, created_on, updated_on, updated_by)
	VALUES
		(:photo_id, :owner_id, :title, :description, :deleted, :created_on, :updated_on, :updated_by)`

	s.log.Printf("%s: %s", "photo.Create", database.Log(query, pht))

	_, err := s.db.NamedExec(query, pht)
	if err != nil {
		return Photo{}, errors.Wrap(err, "inserting photo")
	}

	return pht, nil
}

// Update - update photo
func (s Store) Update(photoID string, up UpdatePhoto) error {
	if err := validate.CheckID(photoID); err != nil {
		return database.ErrInvalidID
	}
	if err := validate.Check(up); err != nil {
		return errors.Wrap(err, "validating data")
	}
	pht, err := s.QueryByID(photoID)
	if err != nil {
		return errors.Wrap(err, "updating photo")
	}

	if up.Title != nil {
		pht.Title = *up.Title
	}
	if up.Description != nil {
		pht.Description = *up.Description
	}
	if up.Deleted != nil {
		pht.Deleted = *up.Deleted
	}
	pht.UpdatedOn = time.Now()

	const query = `
	UPDATE
		photos
	SET
		title = :title,
		description = :description,
		deleted = :deleted,
		updated_on = :updated_on,
		updated_by = :updated_by
	WHERE
		photo_id = :photo_id`

	s.log.Printf("%s: %s", "photo.Update", database.Log(query, pht))

	if _, err = s.db.NamedExec(query, pht); err != nil {
		return errors.Wrapf(err, "updating photo %s", pht.ID)
	}

	return nil
}

// Delete - delete photo from db (and its files)
func (s Store) Delete(photoID string) error {
	if err := validate.CheckID(photoID); err != nil {
		return database.ErrInvalidID
	}
	pht, err := s.QueryByID(photoID)
	if err != nil {
		return errors.Wrap(err, "deleting photo")
	}
	upd := UpdatePhoto{
		Deleted: tests.BoolPointer(true),
	}

	if err := s.Update(pht.ID, upd); err != nil {
		return errors.Wrapf(err, "deleting photo %s", pht.ID)
	}

	const query = "DELETE FROM photo_files WHERE photo_id = ?"

	res, err := s.db.Exec(query, pht.ID)
	if err != nil {
		return errors.Wrapf(err, "deleting photo %s", pht.ID)
	}
	numRows, _ := res.RowsAffected()
	s.log.Println(" -- deleted from photo_files: ", numRows)

	return nil
}

// QueryByID - return given photo
func (s Store) QueryByID(photoID string) (Photo, error) {

	if err := validate.CheckID(photoID); err != nil {
		return Photo{}, database.ErrInvalidID
	}

	data := struct {
		PhotoID string `db:"photo_id"`
	}{
		PhotoID: photoID,
	}
	const query = `
        SELECT photo_id, owner_id, title, description, deleted, created_on, updated_on, updated_by
		FROM photos
		WHERE photo_id = :photo_id`

	s.log.Printf("%s %s", "photo.QueryByID", database.Log(query, data))

	var pht Photo
	if err := database.NamedQueryStruct(s.db, query, data, &pht); err != nil {
		if err == database.ErrNotFound {
			return Photo{}, database.ErrNotFound
		}
		return Photo{}, errors.Wrapf(err, "selecting photo %q", data.PhotoID)
	}

	return pht, nil
}

// QueryByOwnerID - return a list of photos
func (s Store) QueryByOwnerID(ownerID int) ([]Photo, error) {

	data := struct {
		OwnerID int `db:"owner_id"`
	}{
		OwnerID: ownerID,
	}
	const query = `
        SELECT photo_id, owner_id, title, description, created_on, updated_on, updated_by
		FROM photos
		WHERE owner_id = :owner_id`

	s.log.Printf("%s %s", "photo.QueryByOwnerID", database.Log(query, data))

	var photos []Photo
	if err := database.NamedQuerySlice(s.db, query, data, &photos); err != nil {
		/*s.log.Printf("ERR: %s\n", err)
		if err == database.ErrNotFound {
			s.log.Printf("ERR: %s\n", err)
			return nil, database.ErrNotFound
		}*/
		return nil, errors.Wrapf(err, "selecting photos by owner_id %q", data.OwnerID)
	}

	return photos, nil
}

// CreateFile - add new photo into db
func (s Store) CreateFile(npf NewPhotoFile) (PhotoFile, error) {

	if err := validate.Check(npf); err != nil {
		return PhotoFile{}, errors.Wrap(err, "validating data")
	}
	if err := validate.CheckID(npf.PhotoID); err != nil {
		return PhotoFile{}, database.ErrInvalidID
	}

	photoSize, ok := allowedPhotoSizes[npf.Size]
	if !ok {
		return PhotoFile{}, fmt.Errorf("received invalid size: %s", npf.Size)
	}

	now := time.Now()

	phFile := PhotoFile{
		ID:        validate.GenerateID(),
		PhotoID:   npf.PhotoID,
		FilePath:  npf.FilePath,
		Size:      npf.Size,
		Width:     photoSize.w,
		Height:    photoSize.h,
		CreatedOn: now,
		UpdatedOn: now,
		UpdatedBy: npf.UpdatedBy,
	}

	const query = `
	INSERT INTO photo_files
		(file_id, photo_id, filepath, size, w, h, created_on, updated_on, updated_by)
	VALUES
		(:file_id, :photo_id, :filepath, :size, :w, :h, :created_on, :updated_on, :updated_by)`

	_, err := s.db.NamedExec(query, phFile)
	if err != nil {
		s.log.Printf("%s: %s", "photo.Createfile", database.Log(query, phFile))
		return PhotoFile{}, errors.Wrap(err, "inserting photo file")
	}

	return phFile, nil
}

// QueryPhotoFiles - return a list of photos
func (s Store) QueryPhotoFiles(photoID string) ([]PhotoFile, error) {

	data := struct {
		PhotoID string `db:"photo_id"`
	}{
		PhotoID: photoID,
	}
	const query = `
        SELECT file_id, photo_id, filepath, size, w, h, created_on, updated_on, updated_by
		FROM photo_files
		WHERE photo_id = :photo_id`

	s.log.Printf("%s %s", "photo.QueryPhotoFiles", database.Log(query, data))

	var files []PhotoFile
	if err := database.NamedQuerySlice(s.db, query, data, &files); err != nil {
		/*s.log.Printf("ERR: %s\n", err)
		if err == database.ErrNotFound {
			s.log.Printf("ERR: %s\n", err)
			return nil, database.ErrNotFound
		}*/
		return nil, errors.Wrapf(err, "selecting files for photo_id %s", data.PhotoID)
	}

	return files, nil
}
