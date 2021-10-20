package user

import (
	"time"
)

// AuthUser - user
type AuthUser struct {
	ID              int       `db:"user_id" json:"id"`
	Name            string    `db:"name" json:"name"`
	Email           string    `db:"email" json:"email"`
	Pass            []byte    `db:"passw" json:"-"`
	CreatedOn       time.Time `db:"created" json:"date_created"`
	Street          string    `db:"street" json:"street"`
	City            string    `db:"city" json:"city"`
	State           string    `db:"state" json:"state"`
	Zip             string    `db:"zip" json:"zip"`
	Phone           string    `db:"phone" json:"phone"`
	Age             int       `db:"age" json:"age"`
	Gender          string    `db:"gender" json:"gender"`
	Ethnicity       string    `db:"ethnicity" json:"ethnicity"`
	OtherEthnicity  string    `db:"other_ethnicity" json:"other_ethnicity"`
	PermissionLevel int       `db:"permission_level" json:"permission_level"`
}

// NewAuthUser - struct for creating new users
type NewAuthUser struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required"`
	Pass           string `json:"pass" validate:"required"`
	PassConfirm    string `json:"pass_confirm" validate:"eqfield=Pass"`
	Street         string `json:"street" validate:"required"`
	City           string `json:"city" validate:"required"`
	State          string `json:"state" validate:"required"`
	Zip            string `json:"zip" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	Age            int    `json:"age" validate:"required"`
	Gender         string `json:"gender" validate:"required"`
	Ethnicity      string `json:"ethnicity" validate:"required"`
	OtherEthnicity string `json:"other_ethnicity"`
}

// UpdateAuthUser - struct for updating users
type UpdateAuthUser struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Street         string `json:"street" validate:"required"`
	City           string `json:"city" validate:"required"`
	State          string `json:"state" validate:"required"`
	Zip            string `json:"zip" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	Age            int    `json:"age" validate:"required"`
	Gender         string `json:"gender" validate:"required"`
	Ethnicity      string `json:"ethnicity" validate:"required"`
	OtherEthnicity string `json:"other_ethnicity"`
}

// UpdateAuthUserPass - struct for updating users' passwords
type UpdateAuthUserPass struct {
	Pass        string `json:"pass" validate:"required"`
	PassConfirm string `json:"pass_confirm" validate:"eqfield=Pass"`
}

//ResetPasswordEmail - password reset type
type ResetPasswordEmail struct {
	ResetID   string    `db:"reset_id" json:"reset_id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Active    bool      `db:"active" json:"active"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
	UpdatedOn time.Time `db:"updated_on" json:"updated_on"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

//NewResetPasswordEmail - struct for creating new password resets
type NewResetPasswordEmail struct {
	UserID    int    `db:"user_id" json:"user_id"`
	UpdatedBy string `json:"updated_by" validate:"required"`
}

//ExpireResetPasswordEmail - struct for expiring password resets
type ExpireResetPasswordEmail struct {
	ResetID string `json:"reset_id" validate:"required"`
}
