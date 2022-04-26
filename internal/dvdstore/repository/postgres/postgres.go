package repository

import (
	"database/sql"
	"fmt"
)

var (
	ErrCustomerNotFound      = fmt.Errorf("Customer not found")
	ErrOrderNotFound         = fmt.Errorf("Order not found")
	ErrProductOutOfInventory = fmt.Errorf("Product out of inventory")
	ErrProductNotFound       = fmt.Errorf("Product not found")
)

// pgRepo implements PostgresRepo interface
type pgRepo struct {
	db *sql.DB
}

// NewPgRepo is a pgRepo constructor
func NewPgRepo(db *sql.DB) *pgRepo {
	return &pgRepo{db: db}
}
