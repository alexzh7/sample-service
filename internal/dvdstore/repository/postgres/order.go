package repository

import (
	"fmt"
	"sort"
	"time"

	"github.com/alexzh7/sample-service/internal/models"
	"github.com/lib/pq"
)

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
		return nil, models.ErrOrderNotFound
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
		return nil, models.ErrOrderNotFound
	}

	return orders, nil
}

// AddOrder creates order for customerId with provided products. Passed products must have id and
// quantity fields filled. Returns order and errors: ErrOutOfInventory if product is out of inventory,
// ErrProductNotFound if product was not found
// TODO: Add validation on product.Id, product.Quantity > 0
func (p *pgRepo) AddOrder(customerId int, products []*models.Product) (*models.Order, error) {
	// Helper func
	fail := func(errString string, err error) (*models.Order, error) {
		return nil, fmt.Errorf("AddOrder "+errString+": %v ", err)
	}

	// Retrieve product ids for query
	productIds := make([]int, 0)
	for _, p := range products {
		productIds = append(productIds, p.Id)
	}

	selectQuery := `
	SELECT i.prod_id, i.quan_in_stock, p.price, p.title
	FROM inventory i INNER JOIN products p
	ON i.prod_id = p.prod_id
	WHERE i.prod_id = ANY($1)`

	tx, err := p.db.Begin()
	if err != nil {
		return fail("tx.Begin", err)
	}
	defer tx.Rollback()

	// Check products existence and their quantity in stock
	rows, err := tx.Query(selectQuery, pq.Array(productIds))
	if err != nil {
		return fail("SELECT tx.Query", err)
	}
	defer rows.Close()

	prodsInStock := make([]*models.Product, 0)
	for rows.Next() {
		prod := models.Product{}
		if err = rows.Scan(&prod.Id, &prod.Quantity, &prod.Price, &prod.Title); err != nil {
			return fail("SELECT inventory rows.Scan", err)
		}
		prodsInStock = append(prodsInStock, &prod)
	}

	// Check existence
	if len(products) != len(prodsInStock) {
		return nil, models.ErrProductNotFound
	}
	// Check quantity, add products info
	var tax, net, total float64
	sort.Sort(models.SortById(products))
	for i, p := range products {
		if p.Quantity > prodsInStock[i].Quantity {
			return nil, models.ErrProductOutOfInventory
		}
		p.Price = prodsInStock[i].Price
		p.Title = prodsInStock[i].Title
		net += p.Price * float64(p.Quantity)
	}
	tax = net * 0.1
	total = net + tax

	// Update quantity
	// TODO: optimize for one query
	stmt, err := tx.Prepare("UPDATE inventory SET quan_in_stock = quan_in_stock - $1 WHERE prod_id = $2")
	if err != nil {
		return fail("UPDATE inventory tx.Prepare", err)
	}
	defer stmt.Close()
	for _, p := range products {
		if _, err := stmt.Exec(p.Quantity, p.Id); err != nil {
			return fail("UPDATE inventory tx.Exec", err)
		}
	}

	// Insert order
	insertQuery := `
	INSERT INTO orders (orderdate, customerid, netamount, tax, totalamount) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING orderid
	`
	ord := &models.Order{
		Date:        time.Now().UTC(),
		NetAmount:   net,
		Tax:         tax,
		TotalAmount: total,
		Products:    products,
	}
	if err = tx.QueryRow(insertQuery, ord.Date, customerId, ord.NetAmount, ord.Tax, ord.TotalAmount).
		Scan(&ord.Id); err != nil {
		return fail("INSERT orders tx.QueryRow", err)
	}

	if err = tx.Commit(); err != nil {
		return fail("INSERT orders tx.Commit", err)
	}

	return ord, nil
}

// DeleteOrder deletes order by given order id
func (p *pgRepo) DeleteOrder(orderId int) error {
	_, err := p.db.Exec("DELETE FROM orders WHERE orderid=$1", orderId)
	if err != nil {
		return fmt.Errorf("DeleteOrder sql.Exec: %v", err)
	}
	return nil
}
