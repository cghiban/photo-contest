package contest

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

// Create - add new contest into db
func (s Store) Create(nc NewContest) (Contest, error) {

	if err := validate.Check(nc); err != nil {
		return Contest{}, errors.Wrap(err, "validating data")
	}

	now := time.Now().Truncate(time.Second)

	slug := nc.Title

	ctst := Contest{
		Title:       nc.Title,
		Slug:        slug,
		Description: nc.Description,
		StartDate:   nc.StartDate,
		EndDate:     nc.EndDate,
		CreatedOn:   now,
		UpdatedOn:   now,
		UpdatedBy:   nc.UpdatedBy,
	}

	const query = `
	INSERT INTO contests
		(slug, title, description, start_date, end_date, created_on, updated_on, updated_by)
	VALUES
		(:slug, :title, :description, :start_date, :end_date, :created_on, :updated_on, :updated_by)`

	s.log.Printf("%s: %s", "contest.Create", database.Log(query, ctst))

	res, err := s.db.NamedExec(query, ctst)
	if err != nil {
		return Contest{}, errors.Wrap(err, "inserting contest")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Contest{}, err
	}
	ctst.ID = int(id)

	return ctst, nil
}

// CreateContestPhoto - add new contest photo
func (s Store) CreateContestPhoto(ncp NewContestPhoto) (ContestPhoto, error) {

	if err := validate.Check(ncp); err != nil {
		return ContestPhoto{}, errors.Wrap(err, "validating data")
	}

	now := time.Now().Truncate(time.Second)

	cPhoto := ContestPhoto{
		ContestID: ncp.ContestID,
		PhotoID:   ncp.PhotoID,
		Status:    ncp.Status,
		CreatedOn: now,
		UpdatedOn: now,
		UpdatedBy: ncp.UpdatedBy,
	}

	const query = `
	INSERT INTO contest_photos
		(contest_id, photo_id, status, created_on, updated_on, updated_by)
	VALUES
		(:contest_id, :photo_id, :status, :created_on, :updated_on, :updated_by)`

	s.log.Printf("%s: %s", "contest.Create", database.Log(query, cPhoto))

	res, err := s.db.NamedExec(query, cPhoto)
	if err != nil {
		return ContestPhoto{}, errors.Wrap(err, "inserting contest")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return ContestPhoto{}, err
	}
	cPhoto.ID = int(id)

	return cPhoto, nil
}
