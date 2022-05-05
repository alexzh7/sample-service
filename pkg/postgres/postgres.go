package postgres

import (
	"database/sql"
	"fmt"

	"github.com/alexzh7/sample-service/config"
)

// NewPostgresConn returns new connection to postgresql from passed config params
func NewPostgresConn(c *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
