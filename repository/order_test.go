package repository

import (
	"testing"
	"time"

	"github.com/alexzh7/sample-service/models"
	"github.com/stretchr/testify/assert"
)

func TestGetOrder(t *testing.T) {
	o := &models.Order{
		Id:          1,
		Date:        time.Now().UTC(),
		NetAmount:   100.00,
		Tax:         20.00,
		TotalAmount: 120.00,
		Products: []*models.Product{
			{Id: 1, Title: "Interstellar", Price: 80.00, Quantity: 60},
			{Id: 2, Title: "John Wick", Price: 100.00, Quantity: 230},
			{Id: 3, Title: "Inception", Price: 120.00, Quantity: 400},
		},
	}
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

func TestGetCustomerOrders(t *testing.T) {
	orders := []*models.Order{
		{
			Id:          1,
			Date:        time.Now().UTC(),
			NetAmount:   100.00,
			Tax:         20.00,
			TotalAmount: 120.00,
			Products: []*models.Product{
				{Id: 1, Title: "Interstellar", Price: 80.00, Quantity: 60},
				{Id: 2, Title: "John Wick", Price: 100.00, Quantity: 230},
				{Id: 3, Title: "Inception", Price: 120.00, Quantity: 400},
			},
		},
		{
			Id:          2,
			Date:        time.Now().UTC(),
			NetAmount:   200.00,
			Tax:         30.00,
			TotalAmount: 230.00,
			Products: []*models.Product{
				{Id: 55, Title: "Marvel", Price: 90.00, Quantity: 12},
				{Id: 78, Title: "Movie", Price: 60.00, Quantity: 5},
			},
		},
	}

	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{"orderid", "orderdate", "netamount", "tax", "totalamount",
		"prod_id", "title", "price", "quantity"})

	for _, o := range orders {
		for _, p := range o.Products {
			rows.AddRow(o.Id, o.Date, o.NetAmount, o.Tax, o.TotalAmount,
				p.Id, p.Title, p.Price, p.Quantity)
		}
	}

	var customerId = 5
	mock.ExpectQuery("SELECT (.+)").WithArgs(customerId).WillReturnRows(rows)

	ords, err := repo.GetCustomerOrders(customerId)
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(orders, ords) {
		t.Error(NotEqualErr(orders, ords))
	}
}

func TestGetCustomerOrdersNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := &pgRepo{db}

	rows := mock.NewRows([]string{})
	mock.ExpectQuery("SELECT (.+)").WithArgs(10).WillReturnRows(rows)

	order, err := repo.GetCustomerOrders(10)
	assert.ErrorIs(t, err, ErrOrderNotFound)
	assert.Nil(t, order)
}

// проверить на tx.rollback при ErrOutOfInventory и ErrProductNotFound
