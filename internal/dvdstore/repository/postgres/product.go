package repository

import (
	"database/sql"
	"fmt"

	"github.com/alexzh7/sample-service/internal/models"
)

// GetAllProducts returns slice of all products limited by limit
func (p *pgRepo) GetAllProducts(limit int) ([]*models.Product, error) {
	rows, err := p.db.Query(sqlGetAllProducts, limit)
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts sql.Query: %v", err)
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
	for rows.Next() {
		prod := models.Product{}
		if err := rows.Scan(&prod.Id, &prod.Title, &prod.Price, &prod.Quantity); err != nil {
			return nil, fmt.Errorf("GetAllProducts rows.Scan: %v", err)
		}
		products = append(products, &prod)
	}
	if err = rows.Err(); err != nil {
		return products, fmt.Errorf("GetAllProducts rows.Next: %v", err)
	}

	return products, nil
}

// GetProduct returns single product by given id and EntityError if product wasn't found
func (p *pgRepo) GetProduct(productId int) (*models.Product, error) {
	prod := models.Product{}

	err := p.db.QueryRow(sqlGetProduct, productId).
		Scan(&prod.Id, &prod.Title, &prod.Price, &prod.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound("product", productId)
		}
		return nil, fmt.Errorf("GetProduct sql.QueryRow: %v", err)
	}
	return &prod, nil
}

// AddProduct adds a product returning id
func (p *pgRepo) AddProduct(prod *models.Product) (productId int, err error) {
	// Helper func
	fail := func(errSring string, err error) (int, error) {
		return 0, fmt.Errorf("AddProduct "+errSring+": %v", err)
	}

	tx, err := p.db.Begin()
	if err != nil {
		return fail("tx.Begin", err)
	}
	defer tx.Rollback()

	// Insert new product
	if err = tx.QueryRow(sqlAddProduct, prod.Title, prod.Price).Scan(&productId); err != nil {
		return fail("tx.Exec on products", err)
	}

	// Insert quantity
	if _, err = tx.Exec("INSERT INTO inventory (prod_id, quan_in_stock, sales) VALUES ($1, $2, -1)",
		productId, prod.Quantity); err != nil {
		return fail("tx.Exec on inventory", err)
	}

	if err = tx.Commit(); err != nil {
		return fail("tx.Commit", err)
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
