package photo

import "time"

// Photo - photo type
type Photo struct {
	ID          string    `db:"photo_id" json:"id"`
	OwnerID     int       `db:"owner_id" json:"owner_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Deleted     bool      `db:"deleted" json:"deleted"`
	CreatedOn   time.Time `db:"created_on" json:"created_on"`
	UpdatedOn   time.Time `db:"updated_on" json:"updated_on"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}

//NewPhoto - type fore creating a new photo
type NewPhoto struct {
	OwnerID     int    `json:"owner_id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	UpdatedBy   string `json:"updated_by" validate:"required"`
}

//UpdatePhoto - type fore creating a new photo
type UpdatePhoto struct {
	Title       *string `db:"title"  validate:"omitempty"`
	Description *string `db:"description"  validate:"omitempty"`
	Deleted     *bool   `db:"deleted" validate:"omitempty"`
	UpdatedBy   *string `db:"updated_by", validate:"required"`
}

type fileSize struct {
	w uint16
	h uint16
}

// PhotoFileSize - map with allowed photp_file sizes
type PhotoFileSize map[string]fileSize

var allowedPhotoSizes = PhotoFileSize{
	"thumb":  fileSize{200, 200},
	"small":  fileSize{400, 400},
	"medium": fileSize{800, 800},
	"large":  fileSize{1200, 1200},
}

// PhotoFile - photo_file type
type PhotoFile struct {
	ID        string    `db:"file_id" json:"id"`
	PhotoID   string    `db:"photo_id" json:"owner_id"`
	FilePath  string    `db:"filepath" json:"filepath"`
	Size      string    `db:"size" json:"size"`
	Width     uint16    `db:"w" json:"width"`
	Height    uint16    `db:"h" json:"height"`
	Deleted   bool      `db:"deleted" json:"deleted"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
	UpdatedOn time.Time `db:"updated_on" json:"updated_on"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

// NewPhotoFile - type used for creating a new photo_file
type NewPhotoFile struct {
	PhotoID  string `db:"owner_id" json:"owner_id"`
	FilePath string `db:"filepath" json:"filepath"`
	Size     string `db:"size" json:"size"`
	//Width     uint16    `db:"w" json:"width"`
	//Height    uint16    `db:"h" json:"height"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

// UpdatePhotoFile - type for updating a photo_file
type UpdatePhotoFile struct {
	Deleted   bool      `db:"deleted" json:"deleted"`
	UpdatedOn time.Time `db:"updated_on" json:"updated_on"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
