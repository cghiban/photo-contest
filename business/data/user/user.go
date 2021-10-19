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
		Name:           nu.Name,
		Email:          nu.Email,
		Pass:           hash,
		CreatedOn:      time.Now(),
		Street:         nu.Street,
		City:           nu.City,
		State:          nu.State,
		Zip:            nu.Zip,
		Phone:          nu.Phone,
		Age:            nu.Age,
		Gender:         nu.Gender,
		Ethnicity:      nu.Ethnicity,
		OtherEthnicity: nu.OtherEthnicity,
	}

	const query = `
	INSERT INTO auth_user 
		(email, name, passw, created, street, city, state, zip, phone, age, gender, ethnicity, other_ethnicity)
	VALUES 
		(:email, :name, :passw, :created, :street, :city, :state, :zip, :phone, :age, :gender, :ethnicity, :other_ethnicity)`

	s.log.Printf("%s: %s", "user.Create", database.Log(query, usr))

	res, err := s.db.NamedExec(query, usr)
	if err != nil {
		return AuthUser{}, errors.Wrap(err, "inserting user")
	}

	//rowNum, _ := res.RowsAffected()
	//s.log.Println(" -- added to auth_user: ", rowNum)

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
			user_id, name, email, passw, created, street, city, state, zip, phone, age, gender, ethnicity, other_ethnicity
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
        SELECT user_id, name, email, passw, created, street, city, state, zip, phone, age, gender, ethnicity, other_ethnicity
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
		ID:             user_id,
		Name:           uu.Name,
		Email:          uu.Email,
		Street:         uu.Street,
		City:           uu.City,
		State:          uu.State,
		Zip:            uu.Zip,
		Phone:          uu.Phone,
		Age:            uu.Age,
		Gender:         uu.Gender,
		Ethnicity:      uu.Ethnicity,
		OtherEthnicity: uu.OtherEthnicity,
	}

	const query = `
	UPDATE auth_user SET email = :email, name = :name, street = :street, city = :city, state = :state, zip = :zip, phone = :phone, age = :age, gender = :gender, ethnicity = :ethnicity, other_ethnicity = :other_ethnicity WHERE user_id = :user_id`

	s.log.Printf("%s: %s", "user.Update", database.Log(query, usr))

	_, err := s.db.NamedExec(query, usr)
	if err != nil {
		return AuthUser{}, errors.Wrap(err, "updating user")
	}

	// res, err := s.db.NamedExec(...)
	//rowNum, _ := res.RowsAffected()
	//s.log.Println(" -- updated auth_user: ", rowNum)

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

	_, err = s.db.NamedExec(query, usr)
	if err != nil {
		return AuthUser{}, errors.Wrap(err, "updating user password")
	}

	//rowNum, _ := res.RowsAffected()
	//s.log.Println(" -- updated auth_user password: ", rowNum)

	return usr, nil

}

// Create - create a password reset entry
func (s Store) CreatePasswordReset(nr NewResetPasswordEmail) (ResetPasswordEmail, error) {

	if err := validate.Check(nr); err != nil {
		return ResetPasswordEmail{}, errors.Wrap(err, "validating data")
	}
	now := time.Now().Truncate(time.Second)

	rpe := ResetPasswordEmail{
		ResetID:   validate.GenerateID(),
		UserID:    nr.UserID,
		CreatedOn: now,
		UpdatedOn: now,
		UpdatedBy: nr.UpdatedBy,
	}

	const query = `
	INSERT INTO reset_password_email 
		(reset_id, user_id, created_on, updated_on, updated_by)
	VALUES 
		(:reset_id, :user_id, :created_on, :updated_on, :updated_by)`

	s.log.Printf("%s: %s", "user.CreatePasswordReset", database.Log(query, rpe))

	_, err := s.db.NamedExec(query, rpe)
	if err != nil {
		return ResetPasswordEmail{}, errors.Wrap(err, "inserting reset_password_email")
	}

	return rpe, nil
}

// Create - create a password reset entry
func (s Store) ExpirePasswordReset(er ExpireResetPasswordEmail) (ResetPasswordEmail, error) {

	if err := validate.Check(er); err != nil {
		return ResetPasswordEmail{}, errors.Wrap(err, "validating data")
	}

	now := time.Now().Truncate(time.Second)
	rpe := ResetPasswordEmail{
		ResetID:   er.ResetID,
		UpdatedOn: now,
	}

	const query = `
	UPDATE reset_password_email SET active = 0, updated_on = :updated_on
		WHERE reset_id = :reset_id`

	s.log.Printf("%s: %s", "user.ExpirePasswordReset", database.Log(query, rpe))

	_, err := s.db.NamedExec(query, rpe)
	if err != nil {
		return ResetPasswordEmail{}, errors.Wrap(err, "expiring password reset")
	}
	return rpe, nil
}

func (s Store) QueryPasswordResetByID(reset_id string) (ResetPasswordEmail, error) {
	data := struct {
		ResetID string `db:"reset_id"`
	}{
		ResetID: reset_id,
	}
	const query = `
	SELECT reset_id, user_id, active, created_on, updated_on, updated_by
	FROM reset_password_email
	WHERE reset_id = :reset_id AND active = 1 AND created_on >= datetime('now', '-24 hours')`

	var pReset ResetPasswordEmail
	if err := database.NamedQueryStruct(s.db, query, data, &pReset); err != nil {
		if err == database.ErrNotFound {
			return ResetPasswordEmail{}, database.ErrNotFound
		}
		return ResetPasswordEmail{}, errors.Wrapf(err, "selecting password reset %q", data.ResetID)
	}

	return pReset, nil

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
