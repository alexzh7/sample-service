package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexzh7/sample-service/internal/models"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetOrder(t *testing.T) {
	o := &models.Order{
		Id:          1,
		Date:        time.Now().UTC(),
		NetAmount:   100.00,
		Tax:         20.00,
		TotalAmount: 120.00,
		Products:    mockProducts,
	}
	db, mock := NewMock()
	defer db.Close()

	rows := mock.NewRows([]string{"orderid", "orderdate", "netamount", "tax", "totalamount",
		"prod_id", "title", "price", "quantity"})
	for _, v := range o.Products {
		rows.AddRow(o.Id, o.Date, o.NetAmount, o.Tax, o.TotalAmount,
			v.Id, v.Title, v.Price, v.Quantity)
	}

	mock.ExpectQuery("SELECT (.+)").WithArgs(o.Id).WillReturnRows(rows)

	repo := &pgRepo{db}
	order, err := repo.GetOrder(o.Id)
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(o, order) {
		t.Error(NotEqualErr(o, order))
	}
}

func TestGetOrderNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := mock.NewRows([]string{})
	mock.ExpectQuery("SELECT (.+)").WithArgs(10).WillReturnRows(rows)

	repo := &pgRepo{db}
	order, err := repo.GetOrder(10)
	var e *models.EntityError
	assert.ErrorAs(t, err, &e)
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
			Products:    mockProducts,
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

	repo := &pgRepo{db}
	ords, err := repo.GetCustomerOrders(customerId)
	assert.NoError(t, err)
	if !assert.ObjectsAreEqual(orders, ords) {
		t.Error(NotEqualErr(orders, ords))
	}
}

func TestGetCustomerOrdersNotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := mock.NewRows([]string{})
	mock.ExpectQuery("SELECT (.+)").WithArgs(10).WillReturnRows(rows)

	repo := &pgRepo{db}
	order, err := repo.GetCustomerOrders(10)
	var e *models.EntityError
	assert.ErrorAs(t, err, &e)
	assert.Nil(t, order)
}

func TestAddOrder(t *testing.T) {
	customerId := 3
	db, mock := NewMock()
	defer db.Close()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"prod_id", "quan_in_stock", "price", "title"})
	productIds := make([]int, 0)
	for _, p := range mockProducts {
		productIds = append(productIds, p.Id)
		rows.AddRow(p.Id, p.Quantity, p.Price, p.Title)
	}
	mock.ExpectQuery("SELECT (.+)").WithArgs(pq.Array(productIds)).
		WillReturnRows(rows)

	stmt := mock.ExpectPrepare("UPDATE (.+)")
	for _, p := range mockProducts {
		stmt.ExpectExec().WithArgs(p.Quantity, p.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	var tax, net, total float64
	for _, p := range mockProducts {
		net += p.Price * float64(p.Quantity)
	}
	tax = net * 0.1
	total = net + tax

	ord := &models.Order{
		Id:          203,
		NetAmount:   net,
		Tax:         tax,
		TotalAmount: total,
		Products:    mockProducts,
	}

	rows = sqlmock.NewRows([]string{"orderid"}).AddRow(ord.Id)
	mock.ExpectQuery("INSERT (.+)").WithArgs(AnyTime{}, customerId, ord.NetAmount,
		ord.Tax, ord.TotalAmount).WillReturnRows(rows)

	stmt = mock.ExpectPrepare("INSERT (.+)")
	for i, p := range ord.Products {
		stmt.ExpectExec().WithArgs(i+1, ord.Id, p.Id, p.Quantity, AnyTime{}).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	mock.ExpectCommit()

	repo := &pgRepo{db}
	order, err := repo.AddOrder(customerId, mockProducts)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	ord.Date = order.Date
	if !assert.ObjectsAreEqual(ord, order) {
		t.Error(NotEqualErr(ord, order))
	}
}

func TestAddOrderProductNotFound(t *testing.T) {
	customerId := 3
	db, mock := NewMock()
	defer db.Close()

	productIds := make([]int, 0)
	for _, p := range mockProducts {
		productIds = append(productIds, p.Id)
	}
	fewProducts := mockProducts[1:2]

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"prod_id", "quan_in_stock", "price", "title"})
	for _, f := range fewProducts {
		rows.AddRow(f.Id, f.Quantity, f.Price, f.Title)
	}
	mock.ExpectQuery("SELECT (.+)").WithArgs(pq.Array(productIds)).
		WillReturnRows(rows)

	mock.ExpectRollback()

	repo := &pgRepo{db}
	order, err := repo.AddOrder(customerId, mockProducts)
	assert.Nil(t, order)
	var entErr *models.EntityError
	assert.ErrorAs(t, err, &entErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddOrderOutOfInventory(t *testing.T) {
	customerId := 3
	db, mock := NewMock()
	defer db.Close()

	quantity := 2

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"prod_id", "quan_in_stock", "price", "title"})
	productIds := make([]int, 0)
	for _, p := range mockProducts {
		productIds = append(productIds, p.Id)
		rows.AddRow(p.Id, quantity, p.Price, p.Title)
	}
	mock.ExpectQuery("SELECT (.+)").WithArgs(pq.Array(productIds)).
		WillReturnRows(rows)

	mock.ExpectRollback()

	repo := &pgRepo{db}
	order, err := repo.AddOrder(customerId, mockProducts)
	assert.Nil(t, order)
	var entErr *models.EntityError
	assert.ErrorAs(t, err, &entErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteOrder(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	orderId := 10
	mock.ExpectExec("DELETE (.+)").WithArgs(orderId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &pgRepo{db}
	assert.NoError(t, repo.DeleteOrder(orderId))
	assert.NoError(t, mock.ExpectationsWereMet())
}
