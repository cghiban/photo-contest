package data

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// AuthUser - user
type AuthUser struct {
	ID        int
	Name      string
	Email     string
	Pass      string
	CreatedOn time.Time
}

func encryptPassword(rawPass string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(rawPass)))
}

// CheckPasswd - validates password for login
func (u AuthUser) CheckPasswd(rawPass string) bool {

	return u.Pass == encryptPassword(rawPass)
}

//DataStore - db operations
type DataStore struct {
	DB *sqlx.DB
	L  *log.Logger
}

// GetSQLiteVersion -
func (ds *DataStore) GetSQLiteVersion() (string, error) {
	query := `SELECT sqlite_version()`

	//var row *sql.Row
	row := ds.DB.QueryRow(query)

	var version string
	err := row.Scan(&version)

	return version, err
}

// CreateUser - add new user into db
func (ds *DataStore) CreateUser(u *AuthUser) error {
	query := `
	INSERT INTO auth_user (email, name, passw, created)
	VALUES (?, ?, ?, datetime('now','localtime'))
	`
	stmt, err := ds.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	u.Pass = encryptPassword(u.Pass)
	//encPass := encryptPassword(u.Pass)

	res, err := stmt.Exec(u.Email, u.Name, u.Pass)
	if err != nil {
		return err
	}
	//rowNum, _ := res.RowsAffected()
	//ds.L.Println(" -- added videos to DB: ", rowNum)

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)

	return nil
}

// GetUserByEmail - retrieves all records or of the given reader passed as args
func (ds *DataStore) GetUserByEmail(email string) (*AuthUser, error) {

	query := `
        SELECT user_id, name, email, passw, created
		FROM auth_user
		WHERE email = ?`
	row := ds.DB.QueryRow(query, email)

	if row.Err() != nil {
		ds.L.Println(row.Err())
		return nil, row.Err()
	}
	//ds.L.Println(query, email)

	//var day, created, duration string
	var userID, created string
	var u AuthUser
	err := row.Scan(&userID, &u.Name, &u.Email, &u.Pass, &created)
	if err != nil {
		ds.L.Println("nope...")
		ds.L.Println("****", err)
		return nil, err
	}
	UserID, _ := strconv.Atoi(userID)
	u.ID = UserID
	t, _ := time.Parse("2006-01-02T15:04:05Z", created)
	u.CreatedOn = t

	return &u, nil
}

// GetUserByID - return given user
func (ds *DataStore) GetUserByID(user_id int) (*AuthUser, error) {

	query := `
        SELECT user_id, name, email, passw, created
		FROM auth_user
		WHERE user_id = ?`
	row := ds.DB.QueryRow(query, user_id)

	if row.Err() != nil {
		ds.L.Println(row.Err())
		return nil, row.Err()
	}

	var userID, created string
	var u AuthUser
	err := row.Scan(&userID, &u.Name, &u.Email, &u.Pass, &created)
	if err != nil {
		ds.L.Println("nope...")
		ds.L.Println("****", err)
		return nil, err
	}
	UserID, _ := strconv.Atoi(userID)
	u.ID = UserID
	t, _ := time.Parse("2006-01-02T15:04:05Z", created)
	u.CreatedOn = t

	return &u, nil
}
