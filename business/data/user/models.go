package user

import (
	"time"
)

// AuthUser - user
type AuthUser struct {
	ID        int       `db:"user_id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Pass      []byte    `db:"passw" json:"-"`
	CreatedOn time.Time `db:"created" json:"date_created"`
}

// NewAuthUser - struct for creating new users
type NewAuthUser struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Pass        string `json:"pass" valdate:"required"`
	PassConfirm string `json:"pass_confirm" validate:"eqfield=Pass"`
}
