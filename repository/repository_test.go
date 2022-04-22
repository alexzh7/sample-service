package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
)

// Helper function for assert.ObjectsAreEqual
var NotEqualErr = func(want, got interface{}) error {
	return fmt.Errorf("objects not equal, expected: %v, got: %v", want, got)
}

// Mock db connection
// TODO: must be part of tests!
func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating mock db: %v", err)
	}
	return db, mock
}
