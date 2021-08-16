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
