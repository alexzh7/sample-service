package repository

import (
	"database/sql"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
)

// Mock db connection
func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating mock db: %v", err)
	}
	return db, mock
}
