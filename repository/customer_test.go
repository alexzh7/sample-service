package repository

import (
	"database/sql"
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

func TestGetCustomers(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{"customerid", "firstname", "lastname", "age"})
	for _, v := range customers {
		rows.AddRow(v.Id, v.FirstName, v.LastName, v.Age)
	}

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	cst, err := repo.GetCustomers(len(customers))
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
		AddRow(c.Id, c.FirstName, c.LastName, c.Age)

	mock.ExpectQuery("SELECT (.+)").WithArgs(c.Id).WillReturnRows(rows)

	cst, err := repo.GetCustomer(c.Id)
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(c, cst) {
		t.Error(NotEqualErr(c, cst))
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

	prep.ExpectExec().WithArgs(c.FirstName, c.LastName, c.Age).
		WillReturnResult(sqlmock.NewResult(lastInsertId, 1))

	id, err := repo.AddCustomer(c)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, lastInsertId, id)
}

func TestDeleteCustomer(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	prep := mock.ExpectPrepare("DELETE (.+)")

	prep.ExpectExec().WithArgs(c.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteCustomer(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
