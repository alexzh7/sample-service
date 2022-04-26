package repository

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProducts(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := mock.NewRows([]string{"prod_id", "title", "price", "quan_in_stock"})
	for _, v := range mockProducts {
		rows.AddRow(v.Id, v.Title, v.Price, v.Quantity)
	}
	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	repo := &pgRepo{db}
	prods, err := repo.GetAllProducts(len(mockProducts))
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(mockProducts, prods) {
		t.Error(NotEqualErr(mockProducts, prods))
	}
}

func TestGetProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := mock.NewRows([]string{"prod_id", "title", "price", "quan_in_stock"}).
		AddRow(mockProduct.Id, mockProduct.Title, mockProduct.Price, mockProduct.Quantity)
	mock.ExpectQuery("SELECT (.+)").WithArgs(mockProduct.Id).WillReturnRows(rows)

	repo := &pgRepo{db}
	pr, err := repo.GetProduct(mockProduct.Id)
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(mockProduct, pr) {
		t.Error(NotEqualErr(mockProduct, pr))
	}
}

func TestGetProductNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	var id = 11
	mock.ExpectQuery("SELECT (.+)").WithArgs(id).WillReturnError(sql.ErrNoRows)

	repo := &pgRepo{db}
	pr, err := repo.GetProduct(id)
	assert.ErrorIs(t, err, ErrProductNotFound)
	assert.Nil(t, pr)
}

func TestAddProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	mock.ExpectBegin()
	var lastInsertId int = 11
	rows := sqlmock.NewRows([]string{"prod_id"}).AddRow(lastInsertId)
	mock.ExpectQuery("INSERT INTO products (.+)").WithArgs(mockProduct.Title, mockProduct.Price).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO inventory (.+)").WithArgs(lastInsertId, mockProduct.Quantity).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &pgRepo{db}
	id, err := repo.AddProduct(mockProduct)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, lastInsertId, id)
}

func TestAddProductRollback(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	mock.ExpectBegin()
	var lastInsertId = 11
	rows := sqlmock.NewRows([]string{"prod_id"}).AddRow(lastInsertId)
	mock.ExpectQuery("INSERT INTO products (.+)").WithArgs(mockProduct.Title, mockProduct.Price).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO inventory (.+)").WithArgs(lastInsertId, mockProduct.Quantity).
		WillReturnError(fmt.Errorf("rollback"))
	mock.ExpectRollback()

	repo := &pgRepo{db}
	_, err := repo.AddProduct(mockProduct)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Error(t, err)
}

func TestDeleteProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM inventory (.+)").WithArgs(mockProduct.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM products (.+)").WithArgs(mockProduct.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := &pgRepo{db}
	err := repo.DeleteProduct(mockProduct.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func TestDeleteProductRollback(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM inventory (.+)").WithArgs(mockProduct.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM products (.+)").WithArgs(mockProduct.Id).
		WillReturnError(fmt.Errorf("rollback"))
	mock.ExpectRollback()

	repo := &pgRepo{db}
	err := repo.DeleteProduct(mockProduct.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Error(t, err)
}
