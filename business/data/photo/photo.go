package photo

import (
	"log"
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

	s.log.Printf("%s: %s", "user.Create", database.Log(query, pht))

	res, err := s.db.NamedExec(query, pht)
	if err != nil {
		return Photo{}, errors.Wrap(err, "inserting photo")
	}

	rowNum, _ := res.RowsAffected()
	s.log.Println(" -- added to photos: ", rowNum)

	return pht, nil
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
        SELECT photo_id, owner_id, title, description, created_on, updated_on, updated_by
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
