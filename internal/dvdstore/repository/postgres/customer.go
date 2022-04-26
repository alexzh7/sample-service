package repository

import (
	"database/sql"
	"fmt"

	"github.com/alexzh7/sample-service/internal/models"
)

// GetAllCustomers returns list of all customers limited by limit
func (p *pgRepo) GetAllCustomers(limit int) ([]*models.Customer, error) {
	rows, err := p.db.Query("SELECT customerid, firstname, lastname, age FROM customers LIMIT $1", limit)
	if err != nil {
		return nil, fmt.Errorf("GetAllCustomers sql.Query: %v", err)
	}
	defer rows.Close()

	customers := make([]*models.Customer, 0)
	for rows.Next() {
		cst := models.Customer{}
		if err := rows.Scan(&cst.Id, &cst.FirstName, &cst.LastName, &cst.Age); err != nil {
			return nil, fmt.Errorf("GetAllCustomers rows.Scan: %v", err)
		}
		customers = append(customers, &cst)
	}
	if err = rows.Err(); err != nil {
		return customers, fmt.Errorf("GetAllCustomers rows.Next: %v", err)
	}

	return customers, nil
}

// GetCustomer returns single customer by given id and ErrCustomerNotFound if customer wasn't found
func (p *pgRepo) GetCustomer(customerId int) (*models.Customer, error) {
	cst := models.Customer{}
	err := p.db.QueryRow("SELECT customerid, firstname, lastname, age FROM customers WHERE customerid=$1",
		customerId).Scan(&cst.Id, &cst.FirstName, &cst.LastName, &cst.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("GetCustomer sql.QueryRow: %v", err)
	}
	return &cst, nil
}

// AddCustomer adds a customer returning id
// TODO: add validation to check firstname, lastname, age
func (p *pgRepo) AddCustomer(cst *models.Customer) (id int, err error) {
	if err = p.db.QueryRow(addCustomerQuery, cst.FirstName, cst.LastName, cst.Age).
		Scan(&id); err != nil {
		return 0, fmt.Errorf("AddCustomer sql.QueryRow: %v", err)
	}
	return id, nil
}

// DeleteCustomer deletes customer with provided id
func (p *pgRepo) DeleteCustomer(customerId int) error {
	_, err := p.db.Exec("DELETE FROM customers WHERE customerid=$1", customerId)
	if err != nil {
		return fmt.Errorf("DeleteCustomer sql.Exec: %v", err)
	}
	return nil
}

// I use only 3 columns from sample database to simplify the project logic
const addCustomerQuery = `
		INSERT INTO customers (
			firstname,
			lastname,
			address1,
			address2,
			city,
			state,
			zip,
			country,
			region,
			email,
			phone,
			creditcardtype,
			creditcard,
			creditcardexpiration,
			username,
			password,
			age,
			income,
			gender
		)
		VALUES (
			$1,
			$2,
			'',
			'',
			'',
			'',
			-1,
			'',
			-1,
			'',
			'',
			-1,
			'',
			'',
			'',
			'',
			$3,
			-1,
			''
		)
		RETURNING customerid
		`
