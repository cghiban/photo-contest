// Package database provides support for access the database.
package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrForbidden             = errors.New("attempted action is not allowed")
)

// Config is the required properties to use the database.
type Config struct {
	Path        string
	Mode        string
	JournalMode string
	Cache       string
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*sqlx.DB, error) {
	mode := "rw"
	if cfg.Mode != "" {
		mode = cfg.Mode
	}

	journalMode := "WAL"
	if cfg.JournalMode != "" {
		journalMode = cfg.JournalMode
	}
	cache := "shared"
	if cfg.JournalMode != "" {
		cache = cfg.JournalMode
	}

	q := make(url.Values)
	q.Set("_journal_mode", journalMode)
	q.Set("mode", mode)
	q.Set("cache", cache)

	u := url.URL{
		//Scheme:   scheme,
		Path:     cfg.Path,
		RawQuery: q.Encode(),
	}
	conn_str := u.String()
	if cfg.Path == ":memory:" {
		conn_str = ":memory:"
	}

	fmt.Println(conn_str)

	db, err := sqlx.Open("sqlite3", conn_str)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {

	// First check we can ping the database.
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity. Running this query forces a
	// round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

// NamedQuerySlice is a helper function for executing queries that return a
// collection of data to be unmarshaled into a slice.
func NamedQuerySlice(db *sqlx.DB, query string, data interface{}, dest interface{}) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("must provide a pointer to a slice")
	}

	rows, err := db.NamedQuery(query, data)
	if err != nil {
		return err
	}

	slice := val.Elem()
	for rows.Next() {
		v := reflect.New(slice.Type().Elem())
		if err := rows.StructScan(v.Interface()); err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, v.Elem()))
	}

	return nil
}

// NamedQueryStruct is a helper function for executing queries that return a
// single value to be unmarshalled into a struct type.
func NamedQueryStruct(db *sqlx.DB, query string, data interface{}, dest interface{}) error {
	rows, err := db.NamedQuery(query, data)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return ErrNotFound
	}

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}

// Log provides a pretty print version of the query and parameters.
func Log(query string, args ...interface{}) string {
	query, params, err := sqlx.Named(query, args)
	if err != nil {
		return err.Error()
	}

	for _, param := range params {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}

	query = strings.Replace(query, "\t", "", -1)
	query = strings.Replace(query, "\n", " ", -1)

	return fmt.Sprintf("[%s]\n", strings.Trim(query, " "))
}

// GetSQLiteVersion -
func GetSQLiteVersion(db *sqlx.DB) (string, error) {
	const query = `SELECT sqlite_version()`

	//var row *sql.Row
	row := db.QueryRow(query)

	var version string
	err := row.Scan(&version)

	return version, err
}
