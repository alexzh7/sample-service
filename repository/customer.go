package repository

import (
	"database/sql"
	"fmt"

	"github.com/alexzh7/sample-service/models"
)

var ErrCustomerNotFound = fmt.Errorf("Customer not found")

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
func (c *customerPgRepo) GetCustomers(limit int) ([]*models.Customer, error) {
	rows, err := c.db.Query("SELECT customerid, firstname, lastname, age FROM customers LIMIT $1", limit)
	if err != nil {
		return nil, fmt.Errorf("GetCustomers sql.Query: %v", err)
	}
	defer rows.Close()

	customers := make([]*models.Customer, 0)

	for rows.Next() {
		cst := models.Customer{}
		if err := rows.Scan(&cst.Id, &cst.FirstName, &cst.LastName, &cst.Age); err != nil {
			return nil, fmt.Errorf("GetCustomers rows.Scan: %v", err)
		}
		customers = append(customers, &cst)
	}
	if err = rows.Err(); err != nil {
		return customers, fmt.Errorf("GetCustomers rows.Next: %v", err)
	}

	return customers, nil
}

//GetCustomer returns single customer by given id and ErrCustomerNotFound if customer wasn't found
func (c *customerPgRepo) GetCustomer(id int) (*models.Customer, error) {
	cst := models.Customer{}
	err := c.db.QueryRow("SELECT customerid, firstname, lastname, age FROM customers WHERE customerid=$1", id).
		Scan(&cst.Id, &cst.FirstName, &cst.LastName, &cst.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("GetCustomer sql.QueryRow: %v", err)
	}
	return &cst, nil
}

//AddCustomer adds a customer returning id
func (c *customerPgRepo) AddCustomer(cst models.Customer) (id int, err error) {
	//stmt, err := db.Prepare("INSERT INTO users(name) VALUES(?)")

	// INSERT INTO customers (
	// 	firstname,
	// 	lastname,
	// 	address1,
	// 	address2,
	// 	city,
	// 	state,
	// 	zip,
	// 	country,
	// 	region,
	// 	email,
	// 	phone,
	// 	creditcardtype,
	// 	creditcard,
	// 	creditcardexpiration,
	// 	username,
	// 	password,
	// 	age,
	// 	income,
	// 	gender
	//   )
	// VALUES (
	// 	'Test',
	// 	'Test',
	// 	'',
	// 	'',
	// 	'',
	// 	'',
	// 	-1,
	// 	'',
	// 	-1,
	// 	'',
	// 	'',
	// 	-1,
	// 	'',
	// 	'',
	// 	'',
	// 	'',
	// 	26,
	// 	-1,
	// 	''
	//   );

	return id, err
}
