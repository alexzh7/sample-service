package repository

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexzh7/sample-service/models"
	"github.com/stretchr/testify/assert"
)

var products = []*models.Product{
	{Id: 1, Title: "Interstellar", Price: 80.00, Quantity: 60},
	{Id: 2, Title: "John Wick", Price: 100.00, Quantity: 230},
	{Id: 3, Title: "Inception", Price: 120.00, Quantity: 400},
}

var p = &models.Product{Id: 1, Title: "Interstellar", Price: 80.00, Quantity: 60}

func TestGetProducts(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{"prod_id", "title", "price", "quan_in_stock"})
	for _, v := range products {
		rows.AddRow(v.Id, v.Title, v.Price, v.Quantity)
	}

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	prods, err := repo.GetProducts(len(products))
	assert.NoError(t, err)
	assert.Len(t, prods, len(products))

	for k, v := range prods {
		assert.Equal(t, products[k], v)
	}
}

func TestGetProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{"prod_id", "title", "price", "quan_in_stock"}).
		AddRow(p.Id, p.Title, p.Price, p.Quantity)

	mock.ExpectQuery("SELECT (.+)").WithArgs(p.Id).WillReturnRows(rows)

	prod, err := repo.GetProduct(p.Id)
	assert.NoError(t, err)
	assert.Equal(t, p, prod)
}

func TestGetProductNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	// Not existing id
	var id = 11
	repo := &pgRepo{db}

	mock.ExpectQuery("SELECT (.+)").WithArgs(id).WillReturnError(sql.ErrNoRows)

	prod, err := repo.GetProduct(id)
	assert.ErrorIs(t, err, ErrProductNotFound)
	assert.Nil(t, prod)
}

func TestAddProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	var lastInsertId int64 = 11
	repo := &pgRepo{db}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO products (.+)").WithArgs(p.Title, p.Price).
		WillReturnResult(sqlmock.NewResult(lastInsertId, 1))

	mock.ExpectExec("INSERT INTO inventory (.+)").WithArgs(lastInsertId, p.Quantity).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	id, err := repo.AddProduct(p)
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
	mock.ExpectExec("INSERT INTO products (.+)").WithArgs(p.Title, p.Price).
		WillReturnResult(sqlmock.NewResult(lastInsertId, 1))

	mock.ExpectExec("INSERT INTO inventory (.+)").WithArgs(lastInsertId, p.Quantity).
		WillReturnError(fmt.Errorf("rollback"))
	mock.ExpectRollback()

	_, err := repo.AddProduct(p)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Error(t, err)
}

func TestDeleteProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM inventory (.+)").WithArgs(p.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM products (.+)").WithArgs(p.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteProduct(p.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func TestDeleteProductRollback(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM inventory (.+)").WithArgs(p.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM products (.+)").WithArgs(p.Id).
		WillReturnError(fmt.Errorf("rollback"))
	mock.ExpectRollback()

	err := repo.DeleteProduct(p.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Error(t, err)
}
