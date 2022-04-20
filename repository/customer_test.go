package repository

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexzh7/sample-service/models"
	"github.com/stretchr/testify/assert"
)

var customers = []*models.Customer{
	{Id: 1, FirstName: "John", LastName: "Doe", Age: 40},
	{Id: 2, FirstName: "Tony", LastName: "Stark", Age: 33},
	{Id: 3, FirstName: "Alex", LastName: "Zhuravlev", Age: 26},
}

var c = &models.Customer{Id: 1, FirstName: "John", LastName: "Doe", Age: 40}

//Mock db connection
func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating mock db: %v", err)
	}
	return db, mock
}

func TestGetCustomers(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &customerPgRepo{db}

	rows := mock.NewRows([]string{"customerid", "firstname", "lastname", "age"})
	for _, v := range customers {
		rows.AddRow(v.Id, v.FirstName, v.LastName, v.Age)
	}

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	cst, err := repo.GetCustomers(len(customers))
	assert.NoError(t, err)
	assert.Len(t, cst, len(customers))

	for k, v := range cst {
		assert.Equal(t, customers[k], v)
	}
}

func TestGetCustomer(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &customerPgRepo{db}

	rows := mock.NewRows([]string{"customerid", "firstname", "lastname", "age"}).
		AddRow(c.Id, c.FirstName, c.LastName, c.Age)

	mock.ExpectQuery("SELECT (.+)").WithArgs(c.Id).WillReturnRows(rows)

	cst, err := repo.GetCustomer(c.Id)
	assert.NoError(t, err)
	assert.Equal(t, c, cst)
}

func TestGetCustomerNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	//Not existing id
	var id = 11
	repo := &customerPgRepo{db}

	mock.ExpectQuery("SELECT (.+)").WithArgs(id).WillReturnError(sql.ErrNoRows)

	cst, err := repo.GetCustomer(id)
	assert.ErrorIs(t, err, ErrCustomerNotFound)
	assert.Nil(t, cst)
}

func TestAddCustomer(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	var lastInsertId int64 = 11
	repo := &customerPgRepo{db}

	prep := mock.ExpectPrepare("INSERT (.+)")

	prep.ExpectExec().WithArgs(c.FirstName, c.LastName, c.Age).
		WillReturnResult(sqlmock.NewResult(lastInsertId, 1))

	id, err := repo.AddCustomer(c)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, lastInsertId, id)
}
