package models

import (
	"database/sql"

	"go.uber.org/zap"
)

//Customer model
type Customer struct {
	Id        int    `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Age       int    `json:"age,omitempty"`
}

type customerPgRepo struct {
	log *zap.SugaredLogger
	db  *sql.DB
}

//NewCustomerPgRepo customerPgRepo constructor 
func NewCustomerPgRepo(log *zap.SugaredLogger, db *sql.DB) *customerPgRepo {
	return &customerPgRepo{log: log, db: db}
}

//GetCustomers returns list of all customers limited by limit
func (c *customerPgRepo) GetCustomers(limit int) ([]Customer, error) {
	limit = 5
	rows, err := c.db.Query("SELECT customerid, firstname, lastname, age FROM customers LIMIT $1", limit)
	if err != nil {
		c.log.Errorf("GetCustomers db.Query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var cst Customer
		if err := rows.Scan(&cst.Id, &cst.FirstName, &cst.LastName, &cst.Age); err != nil {
			c.log.Errorf("GetCustomers rows.Scan: %v", err)
			return nil, err
		}
		customers = append(customers, cst)
	}
	if err = rows.Err(); err != nil {
		c.log.Errorf("GetCustomers rows.Next: %v", err)
		return customers, err
	}

	return customers, nil
}
