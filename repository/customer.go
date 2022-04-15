package repository

import (
	"database/sql"
	"fmt"

	"github.com/alexzh7/sample-service/models"
)

type customerPgRepo struct {
	db *sql.DB
}

//NewCustomerPgRepo customerPgRepo constructor
func NewCustomerPgRepo(db *sql.DB) *customerPgRepo {
	return &customerPgRepo{db: db}
}

//TODO: Сделать интерфейс для методов, вынести работу с postges в repository/postgres?
//		https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e

//GetCustomers returns list of all customers limited by limit
func (c *customerPgRepo) GetCustomers(limit int) ([]models.Customer, error) {
	rows, err := c.db.Query("SELECT customerid, firstname, lastname, age FROM customers LIMIT $1", limit)
	if err != nil {
		return nil, fmt.Errorf("GetCustomers db.Query: %v", err)
	}
	defer rows.Close()

	var customers []models.Customer

	for rows.Next() {
		cst := models.Customer{}
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
