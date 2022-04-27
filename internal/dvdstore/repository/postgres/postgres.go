package repository

import (
	"database/sql"
)

// pgRepo implements PostgresRepo interface
type pgRepo struct {
	db *sql.DB
}

// NewPgRepo is a pgRepo constructor. Returns error if db is unreachable
func NewPgRepo(db *sql.DB) (*pgRepo, error) {
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &pgRepo{db: db}, nil
}
