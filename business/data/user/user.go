package user

import (
	"log"
	"photo-contest/business/sys/validate"
	"photo-contest/foundation/database"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Store manages the set of API's for user access.
type Store struct {
	log *log.Logger
	db  *sqlx.DB
}

// NewStore constructs a user store for api access.
func NewStore(log *log.Logger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// Create - add new user into db
func (s Store) Create(nu NewAuthUser) (AuthUser, error) {

	if err := validate.Check(nu); err != nil {
		return AuthUser{}, errors.Wrap(err, "validating data")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Pass), bcrypt.DefaultCost)
	if err != nil {
		return AuthUser{}, errors.Wrap(err, "generating password hash")
	}

	usr := AuthUser{
		Name:      nu.Name,
		Email:     nu.Email,
		Pass:      hash,
		CreatedOn: time.Now(),
	}

	const query = `
	INSERT INTO auth_user 
		(email, name, passw, created)
	VALUES 
		(:email, :name, :passw, :created)`

	s.log.Printf("%s: %s", "user.Create", database.Log(query, usr))

	res, err := s.db.NamedExec(query, usr)
	if err != nil {
		return AuthUser{}, errors.Wrap(err, "inserting user")
	}

	rowNum, _ := res.RowsAffected()
	s.log.Println(" -- added to auth_user: ", rowNum)

	id, err := res.LastInsertId()
	if err != nil {
		return AuthUser{}, err
	}
	usr.ID = int(id)

	return usr, nil
}

// QueryByEmail - retrieves user
func (s Store) QueryByEmail(email string) (AuthUser, error) {

	// TODO validate email address

	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}
	const query = `
        SELECT
			user_id, name, email, passw, created
		FROM 
			auth_user
		WHERE email = :email`

	s.log.Printf("%s: %s", "user.QueryByEmail",
		database.Log(query, data),
	)

	var usr AuthUser
	if err := database.NamedQueryStruct(s.db, query, data, &usr); err != nil {
		if err == database.ErrNotFound {
			return AuthUser{}, database.ErrNotFound
		}
		return AuthUser{}, errors.Wrapf(err, "selecting user %q", data.Email)
	}

	return usr, nil
}

// QueryByID - return given user
func (s Store) QueryByID(user_id int) (AuthUser, error) {

	data := struct {
		UserID int `db:"user_id"`
	}{
		UserID: user_id,
	}
	const query = `
        SELECT user_id, name, email, passw, created
		FROM auth_user
		WHERE user_id = :user_id`

	s.log.Printf("%s %s", "user.QueryByID", database.Log(query, data))

	var usr AuthUser
	if err := database.NamedQueryStruct(s.db, query, data, &usr); err != nil {
		if err == database.ErrNotFound {
			return AuthUser{}, database.ErrNotFound
		}
		return AuthUser{}, errors.Wrapf(err, "selecting user %q", data.UserID)
	}

	return usr, nil

}

// Update - update user name and email
func (s Store) Update(user_id int, uu UpdateAuthUser) (AuthUser, error) {

	if err := validate.Check(uu); err != nil {
		return AuthUser{}, errors.Wrap(err, "validating data")
	}

	usr := AuthUser{
		ID:    user_id,
		Name:  uu.Name,
		Email: uu.Email,
	}

	const query = `
	UPDATE auth_user SET email = :email
		AND name = :name WHERE user_id = :user_id`

	s.log.Printf("%s: %s", "user.Update", database.Log(query, usr))

	res, err := s.db.NamedExec(query, usr)
	if err != nil {
		return AuthUser{}, errors.Wrap(err, "updating user")
	}

	rowNum, _ := res.RowsAffected()
	s.log.Println(" -- updated auth_user: ", rowNum)

	return usr, nil

}

// UpdatePass - update user password
func (s Store) UpdatePass(user_id int, up UpdateAuthUserPass) (AuthUser, error) {

	if err := validate.Check(up); err != nil {
		return AuthUser{}, errors.Wrap(err, "validating data")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(up.Pass), bcrypt.DefaultCost)
	if err != nil {
		return AuthUser{}, errors.Wrap(err, "generating password hash")
	}

	usr := AuthUser{
		ID:   user_id,
		Pass: hash,
	}

	const query = `
	UPDATE auth_user SET passw = :passw
		WHERE user_id = :user_id`

	s.log.Printf("%s: %s", "user.UpdatePass", database.Log(query, usr))

	res, err := s.db.NamedExec(query, usr)
	if err != nil {
		return AuthUser{}, errors.Wrap(err, "updating user password")
	}

	rowNum, _ := res.RowsAffected()
	s.log.Println(" -- updated auth_user password: ", rowNum)

	return usr, nil

}

// Authenticate finds a user by their email and verifies their password. On
// success it returns the AuthUser instance
func (s Store) Authenticate(email, password string) (*AuthUser, error) {
	usr, err := s.QueryByEmail(email)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword(usr.Pass, []byte(password)); err != nil {
		return nil, database.ErrAuthenticationFailure
	}
	return &usr, nil
}
