package usecase

import (
	"errors"

	"github.com/alexzh7/sample-service/internal/dvdstore"
	"github.com/alexzh7/sample-service/internal/models"
	"go.uber.org/zap"
)

// dvdstoreUC is a use case for dvdstore. It implements Usecase interface
type dvdstoreUC struct {
	pg  dvdstore.PostgresRepo
	log zap.SugaredLogger
	// TODO: Add validation
}

// NewDvdstoreUC returns new dvd store use case
func NewDvdstoreUC(pg dvdstore.PostgresRepo, log zap.SugaredLogger) *dvdstoreUC {
	return &dvdstoreUC{pg: pg, log: log}
}

// GetCustomers returns list of all customers limited by limit
func (d *dvdstoreUC) GetCustomers(limit int) ([]*models.Customer, error) {
	if limit < 0 {
		return nil, errors.New("GetCustomers: limit must be > 0")
	}
	return d.pg.GetAllCustomers(limit)
}

// GetCustomer returns customer by given id and ErrCustomerNotFound if customer wasn't found
func (d *dvdstoreUC) GetCustomer(customerId int) (*models.Customer, error) {
	return d.pg.GetCustomer(customerId)
}

// AddCustomer adds a customer returning id
// TODO: add validation to check firstname, lastname, age
func (d *dvdstoreUC) AddCustomer(customer *models.Customer) (id int, err error) {
	return d.pg.AddCustomer(customer)
}

// DeleteCustomer deletes customer with provided id
func (d *dvdstoreUC) DeleteCustomer(customerId int) error {
	return d.pg.DeleteCustomer(customerId)
}

// GetProducts returns slice of all products limited by limit
func (d *dvdstoreUC) GetProducts(limit int) ([]*models.Product, error) {
	if limit < 0 {
		return nil, errors.New("GetProducts: limit must be > 0")
	}
	return d.pg.GetAllProducts(limit)
}

// GetProduct returns product by given id and ErrProductNotFound if product wasn't found
func (d *dvdstoreUC) GetProduct(productId int) (*models.Product, error) {
	return d.pg.GetProduct(productId)
}

// AddProduct adds a product returning id
// TODO: add validation on title, price, quantity
func (d *dvdstoreUC) AddProduct(prod *models.Product) (productId int, err error) {
	return d.pg.AddProduct(prod)
}

// DeleteProduct deletes product with provided id
func (d *dvdstoreUC) DeleteProduct(productId int) error {
	return d.pg.DeleteProduct(productId)
}

// GetOrder gets order by order id. Returns ErrOrderNotFound if order was not found
func (d *dvdstoreUC) GetOrder(orderId int) (*models.Order, error) {
	return d.pg.GetOrder(orderId)
}

// GetCustomerOrders gets orders for provided customer id. Returns ErrOrderNotFound if order was not found
func (d *dvdstoreUC) GetCustomerOrders(customerId int) ([]*models.Order, error) {
	return d.pg.GetCustomerOrders(customerId)
}

// AddOrder creates order for customerId with provided products. Passed products must have id and
// quantity fields filled. Returns order and errors: ErrOutOfInventory if product is out of inventory,
// ErrProductNotFound if product was not found
// TODO: Add validation on product.Id, product.Quantity > 0
func (d *dvdstoreUC) AddOrder(customerId int, products []*models.Product) (*models.Order, error) {
	return d.pg.AddOrder(customerId, products)
}

// DeleteOrder deletes order by given order id
func (d *dvdstoreUC) DeleteOrder(orderId int) error {
	return d.pg.DeleteOrder(orderId)
}
