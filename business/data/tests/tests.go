// Package tests contains supporting code for running tests.
package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"photo-contest/business/data/schema"
	"photo-contest/foundation/database"

	"github.com/jmoiron/sqlx"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// NewUnit creates a test in-memory database. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T) (*log.Logger, *sqlx.DB, func()) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	db, err := database.Open(database.Config{
		Path:        ":memory:",
		JournalMode: "MEMORY",
		Mode:        "memory",
		Cache:       "shared",
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		t.Fatalf("Migrating error: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		fmt.Println("!!!!!!!!!!!!!!! in teardown...")
		t.Helper()
		db.Close()

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	log := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	return log, db, teardown
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}
