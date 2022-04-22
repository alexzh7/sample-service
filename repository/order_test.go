package repository

import (
	"testing"
	"time"

	"github.com/alexzh7/sample-service/models"
	"github.com/stretchr/testify/assert"
)

var o = &models.Order{
	Id:          1,
	Date:        time.Now().UTC(),
	NetAmount:   100.00,
	Tax:         20.00,
	TotalAmount: 120.00,
	Products:    products,
}

func TestGetOrder(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{"orderid", "orderdate", "netamount", "tax", "totalamount",
		"prod_id", "title", "price", "quantity"})
	for _, v := range o.Products {
		rows.AddRow(o.Id, o.Date, o.NetAmount, o.Tax, o.TotalAmount,
			v.Id, v.Title, v.Price, v.Quantity)
	}

	mock.ExpectQuery("SELECT (.+)").WithArgs(o.Id).WillReturnRows(rows)

	order, err := repo.GetOrder(o.Id)
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(o, order) {
		t.Error(NotEqualErr(o, order))
	}
}

func TestGetOrderNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{})
	mock.ExpectQuery("SELECT (.+)").WithArgs(10).WillReturnRows(rows)

	order, err := repo.GetOrder(10)
	assert.ErrorIs(t, err, ErrOrderNotFound)
	assert.Nil(t, order)
}
