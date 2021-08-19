package contest

import "time"

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
	ID        int       `db:"id" json:"id"`
	ContestID int       `db:"contest_id" json:"contest_id"`
	PhotoID   string    `db:"photo_id" json:"photo_id"`
	Status    string    `db:"status" json:"status"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
	UpdatedOn time.Time `db:"updated_on" json:"updated_on"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

//NewContestEntry type for adding phtos to a content
type NewContestEntry struct {
	ContestID int    `json:"contest_id" validate:"required"`
	PhotoID   string `json:"photo_id" validate:"required"`
	Status    string `json:"status" validate:"required,oneof=active eliminated withdrawn flagged"`
	UpdatedBy string `json:"updated_by" validate:"required"`
}

//ContestPhoto
//type ContestPhoto struct {
//}
