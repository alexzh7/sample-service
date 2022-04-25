package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexzh7/sample-service/models"
	"github.com/stretchr/testify/assert"
)

// Sample customer
var cust = &models.Customer{Id: 1, FirstName: "John", LastName: "Doe", Age: 40}

func TestGetAllCustomers(t *testing.T) {
	customers := []*models.Customer{
		{Id: 1, FirstName: "John", LastName: "Doe", Age: 40},
		{Id: 2, FirstName: "Tony", LastName: "Stark", Age: 33},
		{Id: 3, FirstName: "Alex", LastName: "Zhuravlev", Age: 26},
	}
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{"customerid", "firstname", "lastname", "age"})
	for _, v := range customers {
		rows.AddRow(v.Id, v.FirstName, v.LastName, v.Age)
	}

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	cst, err := repo.GetAllCustomers(len(customers))
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(customers, cst) {
		t.Error(NotEqualErr(customers, cst))
	}
}

func TestGetCustomer(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{"customerid", "firstname", "lastname", "age"}).
		AddRow(cust.Id, cust.FirstName, cust.LastName, cust.Age)

	mock.ExpectQuery("SELECT (.+)").WithArgs(cust.Id).WillReturnRows(rows)

	cst, err := repo.GetCustomer(cust.Id)
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(cust, cst) {
		t.Error(NotEqualErr(cust, cst))
	}
}

func TestGetCustomerNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	//Not existing id
	var id = 11
	repo := &pgRepo{db}

	mock.ExpectQuery("SELECT (.+)").WithArgs(id).WillReturnError(sql.ErrNoRows)

	cst, err := repo.GetCustomer(id)
	assert.ErrorIs(t, err, ErrCustomerNotFound)
	assert.Nil(t, cst)
}

func TestAddCustomer(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	var lastInsertId int64 = 11
	repo := &pgRepo{db}

	prep := mock.ExpectPrepare("INSERT (.+)")

	prep.ExpectExec().WithArgs(cust.FirstName, cust.LastName, cust.Age).
		WillReturnResult(sqlmock.NewResult(lastInsertId, 1))

	id, err := repo.AddCustomer(cust)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, lastInsertId, id)
}

func TestDeleteCustomer(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	prep := mock.ExpectPrepare("DELETE (.+)")

	prep.ExpectExec().WithArgs(cust.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteCustomer(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
