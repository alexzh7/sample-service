package dvdstore

import (
	"github.com/alexzh7/sample-service/internal/models"
)

// PostgresRepo is used to interact via postgresql
type PostgresRepo interface {
	GetAllCustomers(limit int) ([]*models.Customer, error)
	GetCustomer(customerId int) (*models.Customer, error)
	AddCustomer(customer *models.Customer) (id int, err error)
	DeleteCustomer(customerId int) error

	GetAllProducts(limit int) ([]*models.Product, error)
	GetProduct(productId int) (*models.Product, error)
	AddProduct(prod *models.Product) (productId int, err error)
	DeleteProduct(productId int) error

	GetOrder(orderId int) (*models.Order, error)
	GetCustomerOrders(customerId int) ([]*models.Order, error)
	AddOrder(customerId int, products []*models.Product) (*models.Order, error)
	DeleteOrder(orderId int) error
}

// Usecase is a use case for dvdstore
type Usecase interface {
	GetCustomers(limit int) ([]*models.Customer, error)
	GetCustomer(customerId int) (*models.Customer, error)
	AddCustomer(customer *models.Customer) (id int, err error)
	DeleteCustomer(customerId int) error

	GetProducts(limit int) ([]*models.Product, error)
	GetProduct(productId int) (*models.Product, error)
	AddProduct(prod *models.Product) (productId int, err error)
	DeleteProduct(productId int) error

	GetOrder(orderId int) (*models.Order, error)
	GetCustomerOrders(customerId int) ([]*models.Order, error)
	AddOrder(customerId int, products []*models.Product) (*models.Order, error)
	DeleteOrder(orderId int) error
}
