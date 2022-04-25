package repository

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexzh7/sample-service/models"
	"github.com/stretchr/testify/assert"
)

// Sample product
var prod = &models.Product{Id: 1, Title: "Interstellar", Price: 80.00, Quantity: 60}

func TestGetAllProducts(t *testing.T) {
	products := []*models.Product{
		{Id: 1, Title: "Interstellar", Price: 80.00, Quantity: 60},
		{Id: 2, Title: "John Wick", Price: 100.00, Quantity: 230},
		{Id: 3, Title: "Inception", Price: 120.00, Quantity: 400},
	}
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{"prod_id", "title", "price", "quan_in_stock"})
	for _, v := range products {
		rows.AddRow(v.Id, v.Title, v.Price, v.Quantity)
	}

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	prods, err := repo.GetAllProducts(len(products))
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(products, prods) {
		t.Error(NotEqualErr(products, prods))
	}
}

func TestGetProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{"prod_id", "title", "price", "quan_in_stock"}).
		AddRow(prod.Id, prod.Title, prod.Price, prod.Quantity)

	mock.ExpectQuery("SELECT (.+)").WithArgs(prod.Id).WillReturnRows(rows)

	pr, err := repo.GetProduct(prod.Id)
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(prod, pr) {
		t.Error(NotEqualErr(prod, pr))
	}
}

func TestGetProductNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	// Not existing id
	var id = 11
	repo := &pgRepo{db}

	mock.ExpectQuery("SELECT (.+)").WithArgs(id).WillReturnError(sql.ErrNoRows)

	pr, err := repo.GetProduct(id)
	assert.ErrorIs(t, err, ErrProductNotFound)
	assert.Nil(t, pr)
}

func TestAddProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	var lastInsertId int64 = 11
	repo := &pgRepo{db}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO products (.+)").WithArgs(prod.Title, prod.Price).
		WillReturnResult(sqlmock.NewResult(lastInsertId, 1))

	mock.ExpectExec("INSERT INTO inventory (.+)").WithArgs(lastInsertId, prod.Quantity).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	id, err := repo.AddProduct(prod)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, lastInsertId, id)
}

func TestAddProductRollback(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	var lastInsertId int64 = 11
	repo := &pgRepo{db}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO products (.+)").WithArgs(prod.Title, prod.Price).
		WillReturnResult(sqlmock.NewResult(lastInsertId, 1))

	mock.ExpectExec("INSERT INTO inventory (.+)").WithArgs(lastInsertId, prod.Quantity).
		WillReturnError(fmt.Errorf("rollback"))
	mock.ExpectRollback()

	_, err := repo.AddProduct(prod)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Error(t, err)
}

func TestDeleteProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM inventory (.+)").WithArgs(prod.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM products (.+)").WithArgs(prod.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteProduct(prod.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func TestDeleteProductRollback(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM inventory (.+)").WithArgs(prod.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM products (.+)").WithArgs(prod.Id).
		WillReturnError(fmt.Errorf("rollback"))
	mock.ExpectRollback()

	err := repo.DeleteProduct(prod.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Error(t, err)
}
