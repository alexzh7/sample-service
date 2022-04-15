package models

import (
	"database/sql"
	"fmt"
)

//Customer model
type Customer struct {
	Id        int    `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Age       int    `json:"age,omitempty"`
}

type customerPgRepo struct {
	db *sql.DB
}

//NewCustomerPgRepo customerPgRepo constructor
func NewCustomerPgRepo(db *sql.DB) *customerPgRepo {
	return &customerPgRepo{db: db}
}

//GetCustomers returns list of all customers limited by limit
func (c *customerPgRepo) GetCustomers(limit int) ([]Customer, error) {
	limit = 5
	rows, err := c.db.Query("SELECT customerid, firstname, lastname, age FROM customers LIMIT $1", limit)
	if err != nil {
		return nil, fmt.Errorf("GetCustomers db.Query: %v", err)
	}
	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var cst Customer
		if err := rows.Scan(&cst.Id, &cst.FirstName, &cst.LastName, &cst.Age); err != nil {
			return nil, fmt.Errorf("GetCustomers rows.Scan: %v", err)
		}
		customers = append(customers, cst)
	}
	if err = rows.Err(); err != nil {
		return customers, fmt.Errorf("GetCustomers rows.Next: %v", err)
	}

	return customers, nil
}
