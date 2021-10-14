package contest

import (
	"time"
)

//Contest - contest type
type Contest struct {
	ID          int       `db:"contest_id" json:"contest_id"`
	Title       string    `db:"title" json:"title"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description"`
	StartDate   time.Time `db:"start_date" json:"start_date"`
	EndDate     time.Time `db:"end_date" json:"end_date"`
	CreatedOn   time.Time `db:"created_on" json:"created_on"`
	UpdatedOn   time.Time `db:"updated_on" json:"updated_on"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}

//NewContest - new contest type
type NewContest struct {
	//Slug        string    `db:"slug" json:"slug"`
	Title       string    `db:"title" json:"title" validate:"required"`
	Description string    `db:"description" json:"description" validate:"required"`
	StartDate   time.Time `db:"start_date" json:"start_date" validate:"required"`
	EndDate     time.Time `db:"end_date" json:"end_date" validate:"required,gtfield=StartDate"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}

//ContestEntry - contest entry type
type ContestEntry struct {
	EntryID          int       `db:"entry_id" json:"entry_id"`
	ContestID        int       `db:"contest_id" json:"contest_id"`
	PhotoID          string    `db:"photo_id" json:"photo_id"`
	SubjectName      string    `db:"sname" json:"sname"`
	SubjectAge       string    `db:"sage" json:"sage"`
	SubjectCountry   string    `db:"scountry" json:"scountry"`
	SubjectOrigin    string    `db:"sorigin" json:"sorigin"`
	Location         string    `db:"location" json:"location"`
	SubjectBiography string    `db:"sbiography" json:"sbiography"`
	ReleaseMimeType  string    `db:"release_mime_type" json:"release_mime_type"`
	Status           string    `db:"status" json:"status"`
	CreatedOn        time.Time `db:"created_on" json:"created_on"`
	UpdatedOn        time.Time `db:"updated_on" json:"updated_on"`
	UpdatedBy        string    `db:"updated_by" json:"updated_by"`
}

//NewContestEntry type for adding phtos to a content
type NewContestEntry struct {
	ContestID        int    `json:"contest_id" validate:"required"`
	PhotoID          string `json:"photo_id" validate:"required"`
	Status           string `json:"status" validate:"required,oneof=active eliminated withdrawn flagged"`
	UpdatedBy        string `json:"updated_by" validate:"required"`
	SubjectName      string `json:"sname" validate:"required"`
	SubjectAge       string `json:"sage" validate:"required"`
	SubjectCountry   string `json:"scountry" validate:"required"`
	SubjectOrigin    string `json:"sorigin" validate:"omitempty"`
	Location         string `json:"location" validate:"required"`
	SubjectBiography string `json:"sbiography" validate:"required"`
	ReleaseMimeType  string `json:"release_mime_type" validate:"required"`
}

//ContestPhotoEntry - contest entry with photos and votes totals/scores
type ContestPhotoEntry struct {
	EntryID   int    `db:"entry_id" json:"entry_id"`
	ContestID int    `db:"contest_id" json:"contest_id"`
	PhotoID   string `db:"photo_id" json:"photo_id"`
	Title     string `db:"title" json:"title"`
	Status    string `db:"status" json:"status"`
	Filepath  string `db:"filepath" json:"filepath"`
	Author    string `db:"author" json:"author"`
	Score     int    `db:"score" json:"score"`
}

//ContestPhotoVote
type ContestPhotoVote struct {
	EntryID   int       `db:"v_entry_id" json:"entry_id" validate:"required"`
	ContestID int       `db:"v_contest_id" json:"contest_id" validate:"required"`
	PhotoID   string    `db:"v_photo_id" json:"photo_id" validate:"required"`
	VoterID   int       `db:"v_user_id" json:"voter_id" validate:"required"`
	Score     int       `db:"v_score" json:"score" validate:"min=1,max=2"`
	CreatedOn time.Time `db:"v_created_on" json:"created_on"`
}
