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

// QueryBySlug - return given contest details
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

// CreateContestEntry - add new contest photo
func (s Store) CreateContestEntry(ncp NewContestEntry) (ContestEntry, error) {

	if err := validate.Check(ncp); err != nil {
		return ContestEntry{}, errors.Wrap(err, "validating data")
	}

	now := time.Now().Truncate(time.Second)

	cPhoto := ContestEntry{
		ContestID:        ncp.ContestID,
		PhotoID:          ncp.PhotoID,
		SubjectName:      ncp.SubjectName,
		SubjectAge:       ncp.SubjectAge,
		SubjectCountry:   ncp.SubjectCountry,
		SubjectOrigin:    ncp.SubjectOrigin,
		SubjectBiography: ncp.SubjectBiography,
		Location:         ncp.Location,
		ReleaseMimeType:  ncp.ReleaseMimeType,
		Status:           ncp.Status,
		CreatedOn:        now,
		UpdatedOn:        now,
		UpdatedBy:        ncp.UpdatedBy,
	}

	const query = `
	INSERT INTO contest_entries
		(contest_id, photo_id, sname, sage, scountry, sorigin, sbiography, location, release_mime_type, status, created_on, updated_on, updated_by)
	VALUES
		(:contest_id, :photo_id, :sname, :sage, :scountry, :sorigin, :sbiography, :location, :release_mime_type, :status, :created_on, :updated_on, :updated_by)`

	s.log.Printf("%s: %s", "contest.Create", database.Log(query, cPhoto))

	res, err := s.db.NamedExec(query, cPhoto)
	if err != nil {
		return ContestEntry{}, errors.Wrap(err, "inserting contest")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return ContestEntry{}, err
	}
	cPhoto.EntryID = int(id)

	return cPhoto, nil
}

// QueryContestEntrys - return a list of entries
func (s Store) QueryContestEntries(contestID int) ([]ContestEntry, error) {

	data := struct {
		ContestID int `db:"contest_id"`
	}{
		ContestID: contestID,
	}
	const query = `
	SELECT entry_id, contest_id, photo_id, sname, sage, scountry, sorigin, sbiography, location, release_mime_type, status, created_on, updated_on, updated_by
	FROM contest_entries
	WHERE contest_id = :contest_id`

	s.log.Printf("%s %s", "contest.QueryContestEntries", database.Log(query, data))

	var cEntries []ContestEntry
	if err := database.NamedQuerySlice(s.db, query, data, &cEntries); err != nil {
		/*s.log.Printf("ERR: %s\n", err)
		if err == database.ErrNotFound {
			s.log.Printf("ERR: %s\n", err)
			return nil, database.ErrNotFound
		}*/
		return nil, errors.Wrapf(err, "selecting contest photos %q", data.ContestID)
	}

	return cEntries, nil
}

// CreateContestPhotoVote - add save the votes
func (s Store) CreateContestPhotoVote(v ContestPhotoVote) error {

	if err := validate.Check(v); err != nil {
		return errors.Wrap(err, "validating data")
	}

	const query = `
	INSERT INTO contest_entry_votes
		(v_entry_id, v_contest_id, v_photo_id, v_user_id, v_score, v_created_on)
	VALUES
		(:v_entry_id, :v_contest_id, :v_photo_id, :v_user_id, :v_score, :v_created_on)`

	s.log.Printf("%s: %s", "contest.CreateContestPhotoVote", database.Log(query, v))

	_, err := s.db.NamedExec(query, v)
	if err != nil {
		return errors.Wrap(err, "inserting vote")
	}

	/*id, err := res.LastInsertId()
	if err != nil {
		return ContestEntry{}, err
	}
	cPhoto.EntryID = int(id)*/

	return nil
}

//QueryContestPhotos retrives the entries w/ the votes/scores
func (s Store) QueryContestPhotos(contestID int) ([]ContestPhotoEntry, error) {

	data := struct {
		ContestID int `db:"contest_id"`
	}{
		ContestID: contestID,
	}

	const query = `
	SELECT e.entry_id, e.contest_id, p.photo_id, e.status, p.title, u.name AS author,
	CASE WHEN f.file_id IS NOT NULL THEN f.filepath ELSE '' END AS filepath,
	CASE WHEN v_id IS NOT NULL THEN sum(v_score) ELSE 0 END AS score
	FROM contest_entries e
	JOIN photos p ON e.photo_id = p.photo_id
	JOIN auth_user u ON p.owner_id = u.user_id
	LEFT JOIN photo_files f ON (p.photo_id = f.photo_id AND f.size = 'small')
	LEFT JOIN contest_entry_votes ON e.entry_id = v_entry_id
	WHERE e.contest_id = :contest_id AND p.deleted = false
	GROUP BY e.entry_id, e.contest_id, p.photo_id, e.status, p.title, author, filepath`

	s.log.Printf("%s %s", "contest.QueryContestPhotos", database.Log(query, data))

	var cpEntries []ContestPhotoEntry
	if err := database.NamedQuerySlice(s.db, query, data, &cpEntries); err != nil {
		/*s.log.Printf("ERR: %s\n", err)
		if err == database.ErrNotFound {
			s.log.Printf("ERR: %s\n", err)
			return nil, database.ErrNotFound
		}*/
		return nil, errors.Wrapf(err, "selecting contest photos %q", data.ContestID)
	}

	return cpEntries, nil
}
