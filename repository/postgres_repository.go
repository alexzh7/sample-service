package repository

import "database/sql"

type pgRepo struct {
	db *sql.DB
}

// TODO: Сделать интерфейс для методов, вынести работу с postges в repository/postgres?
//		https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e

// NewPgRepo productPgRepo constructor
func NewPgRepo(db *sql.DB) *pgRepo {
	return &pgRepo{db: db}
}
