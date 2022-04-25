package repository

import (
	"database/sql"
	"fmt"

	"github.com/alexzh7/sample-service/models"
)

var ErrProductNotFound = fmt.Errorf("Product not found")

// GetAllProducts returns list of all products limited by limit
func (p *pgRepo) GetAllProducts(limit int) ([]*models.Product, error) {
	query := `
	SELECT p.prod_id, p.title, p.price, i.quan_in_stock 
	FROM products p INNER JOIN inventory i
	ON p.prod_id = i.prod_id
	LIMIT $1
	`

	rows, err := p.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("GetProducts sql.Query: %v", err)
	}
	defer rows.Close()

	products := make([]*models.Product, 0)

	for rows.Next() {
		prod := models.Product{}
		if err := rows.Scan(&prod.Id, &prod.Title, &prod.Price, &prod.Quantity); err != nil {
			return nil, fmt.Errorf("GetProducts rows.Scan: %v", err)
		}
		products = append(products, &prod)
	}
	if err = rows.Err(); err != nil {
		return products, fmt.Errorf("GetProducts rows.Next: %v", err)
	}

	return products, nil
}

// GetProduct returns single product by given id and ErrProductNotFound if product wasn't found
func (p *pgRepo) GetProduct(productId int) (*models.Product, error) {
	query := `
	SELECT p.prod_id, p.title, p.price, i.quan_in_stock 
	FROM products p INNER JOIN inventory i
	ON p.prod_id = i.prod_id
	WHERE p.prod_id = $1
	`
	prod := models.Product{}

	err := p.db.QueryRow(query, productId).
		Scan(&prod.Id, &prod.Title, &prod.Price, &prod.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("GetProduct sql.QueryRow: %v", err)
	}
	return &prod, nil
}

// AddProduct adds a product returning id
// TODO: add validation
func (p *pgRepo) AddProduct(prod *models.Product) (productId int64, err error) {

	fail := func(err error) (int64, error) {
		return 0, fmt.Errorf("AddProduct %v ", err)
	}

	// I use only 2 columns from sample database to simplify the project logic
	productsQuery := `
	INSERT INTO products (category, title, actor, price, special, common_prod_id)
	VALUES (-1, $1, '', $2, -1, -1)
	`

	// Insert product to products table and quantity to inventory table
	tx, err := p.db.Begin()
	if err != nil {
		return fail(fmt.Errorf("tx.Begin: %v", err))
	}
	defer tx.Rollback()

	// Insert new product
	res, err := tx.Exec(productsQuery, prod.Title, prod.Price)
	if err != nil {
		return fail(fmt.Errorf("tx.Exec on products: %v", err))
	}
	productId, err = res.LastInsertId()
	if err != nil {
		return fail(fmt.Errorf("sql.LastInsertId on products: %v", err))
	}

	// Insert quantity
	_, err = tx.Exec("INSERT INTO inventory (prod_id, quan_in_stock, sales) VALUES ($1, $2, -1)",
		productId, prod.Quantity)
	if err != nil {
		return fail(fmt.Errorf("tx.Exec on inventory: %v", err))
	}

	if err = tx.Commit(); err != nil {
		return fail(fmt.Errorf("tx.Commit: %v", err))
	}

	return productId, nil
}

// DeleteProduct deletes product with provided id
func (p *pgRepo) DeleteProduct(productId int) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("DeleteProduct tx.Begin: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM inventory WHERE prod_id=$1", productId)
	if err != nil {
		return fmt.Errorf("DeleteProduct tx.Exec on inventory: %v", err)
	}

	_, err = tx.Exec("DELETE FROM products WHERE prod_id=$1", productId)
	if err != nil {
		return fmt.Errorf("DeleteProduct tx.Exec on products: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("DeleteProduct tx.Commit: %v", err)
	}

	return nil
}
