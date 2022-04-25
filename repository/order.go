package repository

import (
	"fmt"
	"sort"

	"github.com/alexzh7/sample-service/models"
	"github.com/lib/pq"
)

// TODO: add errors with product id, order id, etc...
var ErrOrderNotFound = fmt.Errorf("Order not found")
var ErrOutOfInventory = fmt.Errorf("Product out of inventory")

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

// AddOrder creates order for customerId with provided products.
// Passed products must have id and quantity fields filled.
// Returns order and errors: ErrOutOfInventory if product is out of inventory,
// ErrProductNotFound if product was not found
// TODO: Add validation on product.Id, product.Quantity > 0
func (p *pgRepo) AddOrder(customerId int, products []*models.Product) (*models.Order, error) {
	// Retrieve product ids from struct, count amounts
	// var tax float64 = 10.00
	var netAmount float64
	productIds := make([]int, 0)

	for _, v := range products {
		netAmount += v.Price
		// totalAmount = netAmount + netAmount/tax
		productIds = append(productIds, v.Id)
	}

	query := `
	SELECT i.prod_id, i.quan_in_stock, p.price
	FROM inventory i INNER JOIN products p
	ON i.prod_id = p.prod_id
	WHERE prod_id = ANY($1)`

	tx, err := p.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("AddOrder tx.Begin: %v", err)
	}
	defer tx.Rollback()

	// Check products existence and their quantity in stock
	rows, err := tx.Query(query, pq.Array(productIds))
	if err != nil {
		return nil, fmt.Errorf("AddOrder tx.Query: %v", err)
	}
	defer rows.Close()

	prodsInStock := make([]*models.Product, 0)

	for rows.Next() {
		prod := models.Product{}
		if err = rows.Scan(&prod.Id, &prod.Quantity, &prod.Price); err != nil {
			return nil, fmt.Errorf("AddOrder select inventory rows.Scan: %v", err)
		}
		prodsInStock = append(prodsInStock, &prod)
	}

	// Check existence
	if len(products) != len(prodsInStock) {
		return nil, ErrProductNotFound
	}
	// Check quantity
	sort.Sort(models.SortById(products))
	for i, p := range products {
		if p.Quantity > prodsInStock[i].Quantity {
			return nil, ErrOutOfInventory
		}
	}

	// Update quantity
	// TODO: optimize for one query
	stmt, err := tx.Prepare("UPDATE inventory SET quan_in_stock = quan_in_stock - $1 WHERE prod_id = $2")
	if err != nil {
		return nil, fmt.Errorf("AddOrder inventory tx.Prepare: %v", err)
	}
	defer stmt.Close()
	for _, v := range products {
		if _, err := stmt.Exec(v.Quantity, v.Id); err != nil {
			return nil, fmt.Errorf("AddOrder update inventory tx.Exec: %v", err)
		}
	}

	// order := &models.Order{
	// 	Date: time.Now().UTC(),
	// }
	// Взять цену у каждого продукта и умножить на кол-во
	// Сложить цены - netamount
	// Посчитать налог - tax, записать сумму в totalamount
	// INSERT

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("AddOrder tx.Commit: %v", err)
	}

	return nil, nil
}

// DeleteOrder deletes order by given order id
func (p *pgRepo) DeleteOrder(orderId int) error {

	// TX that checks quantity for products, returns ErrOutOfInventory if err
	// Then adds order, returns id

	return nil
}
