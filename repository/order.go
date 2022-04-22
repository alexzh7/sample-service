package repository

import (
	"fmt"

	"github.com/alexzh7/sample-service/models"
)

var ErrOrderNotFound = fmt.Errorf("Order not found")

// GetOrder gets order by order id. Returns ErrOrderNotFound if order was not found
func (p *pgRepo) GetOrder(orderId int) (*models.Order, error) {
	query := `
	SELECT t.orderid, t.orderdate, t.netamount, t.tax, t.totalamount,
	t.prod_id, p.title, p.price, t.quantity
	FROM products p INNER JOIN

	(SELECT o.*, ol.prod_id, ol.quantity
	FROM orders o INNER JOIN orderlines ol
	ON o.orderid = ol.orderid) t

	ON p.prod_id = t.prod_id
	WHERE orderid=$1
	`
	rows, err := p.db.Query(query, orderId)
	if err != nil {
		return nil, fmt.Errorf("GetOrder sql.Query: %v", err)
	}
	defer rows.Close()

	ord := models.Order{}
	products := make([]*models.Product, 0)

	for rows.Next() {
		pr := models.Product{}
		if err := rows.Scan(&ord.Id, &ord.Date, &ord.NetAmount, &ord.Tax, &ord.TotalAmount,
			&pr.Id, &pr.Title, &pr.Price, &pr.Quantity); err != nil {
			return nil, fmt.Errorf("GetOrder rows.Scan: %v", err)
		}
		products = append(products, &pr)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetOrder rows.Next: %v", err)
	}

	if len(products) == 0 {
		return nil, ErrOrderNotFound
	}

	ord.Products = products

	return &ord, nil
}

// GetCustomerOrders gets orders for provided customer id. Returns ErrOrderNotFound if order was not found
func (p *pgRepo) GetCustomerOrders(customerId int) ([]*models.Order, error) {
	query := `
	SELECT t.orderid, t.orderdate, t.netamount, t.tax, t.totalamount,
	t.prod_id, p.title, p.price, t.quantity
	FROM products p INNER JOIN

	(SELECT o.*, ol.prod_id, ol.quantity
	FROM orders o INNER JOIN orderlines ol
	ON o.orderid = ol.orderid) t

	ON p.prod_id = t.prod_id
	WHERE customerid=$1
	`
	rows, err := p.db.Query(query, customerId)
	if err != nil {
		return nil, fmt.Errorf("GetCustomerOrders sql.Query: %v", err)
	}
	defer rows.Close()

	orders := make([]*models.Order, 0)
	var i, id int

	for rows.Next() {
		ord := models.Order{}
		pr := models.Product{}
		if err := rows.Scan(&ord.Id, &ord.Date, &ord.NetAmount, &ord.Tax, &ord.TotalAmount,
			&pr.Id, &pr.Title, &pr.Price, &pr.Quantity); err != nil {
			return nil, fmt.Errorf("GetCustomerOrders rows.Scan: %v", err)
		}
		// Separate different orders and populate them with products
		if id != ord.Id {
			orders = append(orders, &ord)
			id = ord.Id
			i++
		}
		orders[i-1].Products = append(orders[i-1].Products, &pr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCustomerOrders rows.Next: %v", err)
	}

	if len(orders) == 0 {
		return nil, ErrOrderNotFound
	}

	return orders, nil
}

// AddOrder creates order returning id
func (p *pgRepo) AddOrder(order *models.Order) (orderId int, err error) {
	return 0, nil
}

// DeleteOrder
