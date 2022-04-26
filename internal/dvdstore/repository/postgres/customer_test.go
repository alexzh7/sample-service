package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexzh7/sample-service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCustomers(t *testing.T) {
	customers := []*models.Customer{
		{Id: 1, FirstName: "John", LastName: "Doe", Age: 40},
		{Id: 2, FirstName: "Tony", LastName: "Stark", Age: 33},
		{Id: 3, FirstName: "Alex", LastName: "Zhuravlev", Age: 26},
	}
	db, mock := NewMock()
	defer db.Close()

	rows := mock.NewRows([]string{"customerid", "firstname", "lastname", "age"})
	for _, v := range customers {
		rows.AddRow(v.Id, v.FirstName, v.LastName, v.Age)
	}
	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	repo := &pgRepo{db}
	cst, err := repo.GetAllCustomers(len(customers))
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(customers, cst) {
		t.Error(NotEqualErr(customers, cst))
	}
}

func TestGetCustomer(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := mock.NewRows([]string{"customerid", "firstname", "lastname", "age"}).
		AddRow(mockCustomer.Id, mockCustomer.FirstName, mockCustomer.LastName, mockCustomer.Age)
	mock.ExpectQuery("SELECT (.+)").WithArgs(mockCustomer.Id).WillReturnRows(rows)

	repo := &pgRepo{db}
	cst, err := repo.GetCustomer(mockCustomer.Id)
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(mockCustomer, cst) {
		t.Error(NotEqualErr(mockCustomer, cst))
	}
}

func TestGetCustomerNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	var id = 11
	mock.ExpectQuery("SELECT (.+)").WithArgs(id).WillReturnError(sql.ErrNoRows)

	repo := &pgRepo{db}
	cst, err := repo.GetCustomer(id)
	assert.ErrorIs(t, err, ErrCustomerNotFound)
	assert.Nil(t, cst)
}

func TestAddCustomer(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	var lastInsertId = 11
	rows := mock.NewRows([]string{"customerid"}).AddRow(lastInsertId)
	mock.ExpectQuery("INSERT (.+)").WithArgs(mockCustomer.FirstName, mockCustomer.LastName, mockCustomer.Age).
		WillReturnRows(rows)

	repo := &pgRepo{db}
	id, err := repo.AddCustomer(mockCustomer)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, lastInsertId, id)
}

func TestDeleteCustomer(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	mock.ExpectExec("DELETE (.+)").WithArgs(mockCustomer.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &pgRepo{db}
	err := repo.DeleteCustomer(mockCustomer.Id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
