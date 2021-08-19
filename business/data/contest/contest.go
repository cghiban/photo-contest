package contest

import (
	"log"
	"photo-contest/business/sys/validate"
	"photo-contest/foundation/database"
	"time"

	"github.com/avelino/slugify"
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

	ctst := Contest{
		Title:       nc.Title,
		Slug:        slugify.Slugify(nc.Title),
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

// QueryByID - return given contest details
func (s Store) QueryByID(contestID int) (Contest, error) {

	data := struct {
		ContestID int `db:"contest_id"`
	}{
		ContestID: contestID,
	}
	const query = `
	SELECT contest_id, slug, title, description, start_date, end_date, created_on, updated_on, updated_by
	FROM contests
	WHERE contest_id = :contest_id`

	s.log.Printf("%s %s", "contest.QueryByID", database.Log(query, data))

	var c Contest
	if err := database.NamedQueryStruct(s.db, query, data, &c); err != nil {
		if err == database.ErrNotFound {
			return Contest{}, database.ErrNotFound
		}
		return Contest{}, errors.Wrapf(err, "selecting contest %q", data.ContestID)
	}

	return c, nil
}

// QueryByID - return given contest details
func (s Store) QueryBySlug(slug string) (Contest, error) {

	data := struct {
		Slug string `db:"slug"`
	}{
		Slug: slug,
	}
	const query = `
	SELECT contest_id, slug, title, description, start_date, end_date, created_on, updated_on, updated_by
	FROM contests
	WHERE slug = :slug`

	s.log.Printf("%s %s", "contest.QueryBySlug", database.Log(query, data))

	var c Contest
	if err := database.NamedQueryStruct(s.db, query, data, &c); err != nil {
		if err == database.ErrNotFound {
			return Contest{}, database.ErrNotFound
		}
		return Contest{}, errors.Wrapf(err, "selecting contest by slug %q", data.Slug)
	}

	return c, nil
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

// QueryContestPhotos - return a list of photos
func (s Store) QueryContestPhotos(contestID int) ([]ContestPhoto, error) {

	data := struct {
		ContestID int `db:"contest_id"`
	}{
		ContestID: contestID,
	}
	const query = `
	SELECT id, contest_id, photo_id, status, created_on, updated_on, updated_by
	FROM contest_photos
	WHERE contest_id = :contest_id`

	s.log.Printf("%s %s", "contest.QueryContestPhotos", database.Log(query, data))

	var cPhotos []ContestPhoto
	if err := database.NamedQuerySlice(s.db, query, data, &cPhotos); err != nil {
		/*s.log.Printf("ERR: %s\n", err)
		if err == database.ErrNotFound {
			s.log.Printf("ERR: %s\n", err)
			return nil, database.ErrNotFound
		}*/
		return nil, errors.Wrapf(err, "selecting contest photos %q", data.ContestID)
	}

	return cPhotos, nil
}
