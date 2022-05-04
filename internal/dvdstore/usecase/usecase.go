package usecase

import (
	"errors"
	"fmt"

	"github.com/alexzh7/sample-service/internal/dvdstore"
	"github.com/alexzh7/sample-service/internal/models"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// dvdstoreUC is a use case for dvdstore. It implements Usecase interface
type dvdstoreUC struct {
	pg       dvdstore.PostgresRepo
	log      *zap.SugaredLogger
	validate *validator.Validate
}

// NewDvdstoreUC returns new dvd store use case
func NewDvdstoreUC(
	pg dvdstore.PostgresRepo,
	log *zap.SugaredLogger,
	vl *validator.Validate,
) *dvdstoreUC {
	return &dvdstoreUC{pg: pg, log: log, validate: vl}
}

// GetCustomers returns list of all customers limited by limit and ErrGeneralDBFail
// if db returned db-specific error. Limit must be > 0
func (d *dvdstoreUC) GetCustomers(limit int) ([]*models.Customer, error) {
	if err := validateVar(limit, "limit"); err != nil {
		d.log.Debugf("GetCustomers validate.Var: %v", err)
		return nil, err
	}

	customers, err := d.pg.GetAllCustomers(limit)
	if err != nil {
		d.log.Error(err)
		return nil, models.ErrGeneralDBFail
	}

	return customers, nil
}

// GetCustomer returns customer by given id, EntityError if customer wasn't found
// and ErrGeneralDBFail if db returned db-specific error
func (d *dvdstoreUC) GetCustomer(customerId int) (*models.Customer, error) {
	if err := validateVar(customerId, "customerId"); err != nil {
		d.log.Debugf("GetCustomer validate.Var: %v", err)
		return nil, err
	}

	customer, err := d.pg.GetCustomer(customerId)
	if err != nil {
		if _, ok := err.(*models.EntityError); ok {
			return nil, err
		}
		d.log.Error(err)
		return nil, models.ErrGeneralDBFail
	}

	return customer, nil
}

// AddCustomer adds a customer returning id and ErrGeneralDBFail if db returned db-specific error
func (d *dvdstoreUC) AddCustomer(customer *models.Customer) (id int, err error) {
	err = d.validate.StructPartial(customer, "FirstName", "LastName", "Age")
	if err != nil {
		d.log.Debugf("AddCustomer validate.StructPartial: %v", err)
		return 0, models.ErrFieldsNotValid("firstname", "lastname", "age")
	}

	id, err = d.pg.AddCustomer(customer)
	if err != nil {
		d.log.Error(err)
		return 0, models.ErrGeneralDBFail
	}

	return id, nil
}

// DeleteCustomer deletes customer with provided id and ErrGeneralDBFail if db returned
// db-specific error
func (d *dvdstoreUC) DeleteCustomer(customerId int) error {
	if err := validateVar(customerId, "customerId"); err != nil {
		d.log.Debugf("DeleteCustomer validate.Var: %v", err)
		return err
	}

	err := d.pg.DeleteCustomer(customerId)
	if err != nil {
		d.log.Error(err)
		return models.ErrGeneralDBFail
	}

	return nil
}

// GetProducts returns slice of all products limited by limit and ErrGeneralDBFail if db
// returned db-specific error
func (d *dvdstoreUC) GetProducts(limit int) ([]*models.Product, error) {
	if err := validateVar(limit, "limit"); err != nil {
		d.log.Debugf("GetProducts validate.Var: %v", err)
		return nil, err
	}

	products, err := d.pg.GetAllProducts(limit)
	if err != nil {
		d.log.Error(err)
		return nil, models.ErrGeneralDBFail
	}

	return products, nil
}

// GetProduct returns product by given id, EntityError if product wasn't found
// and ErrGeneralDBFail if db returned db-specific error
func (d *dvdstoreUC) GetProduct(productId int) (*models.Product, error) {
	if err := validateVar(productId, "productId"); err != nil {
		d.log.Debugf("GetProduct validate.Var: %v", err)
		return nil, err
	}

	product, err := d.pg.GetProduct(productId)
	if err != nil {
		if _, ok := err.(*models.EntityError); ok {
			return nil, err
		}
		d.log.Error(err)
		return nil, models.ErrGeneralDBFail
	}

	return product, nil
}

// AddProduct adds a product returning id and ErrGeneralDBFail if db returned db-specific error
func (d *dvdstoreUC) AddProduct(prod *models.Product) (productId int, err error) {
	err = d.validate.StructPartial(prod, "Title", "Price", "Quantity")
	if err != nil {
		d.log.Debugf("AddProduct validate.StructPartial: %v", err)
		return 0, models.ErrFieldsNotValid("title", "price", "quantity")
	}

	productId, err = d.pg.AddProduct(prod)
	if err != nil {
		d.log.Error(err)
		return 0, models.ErrGeneralDBFail
	}

	return productId, nil
}

// DeleteProduct deletes product with provided id and ErrGeneralDBFail if db returned db-specific error
func (d *dvdstoreUC) DeleteProduct(productId int) error {
	if err := validateVar(productId, "productId"); err != nil {
		d.log.Debugf("DeleteProduct validate.Var: %v", err)
		return err
	}

	err := d.pg.DeleteProduct(productId)
	if err != nil {
		d.log.Error(err)
		return models.ErrGeneralDBFail
	}

	return nil
}

// GetOrder gets order by order id. Returns EntityError if order was not found
// and ErrGeneralDBFail if db returned db-specific error
func (d *dvdstoreUC) GetOrder(orderId int) (*models.Order, error) {
	if err := validateVar(orderId, "orderId"); err != nil {
		d.log.Debugf("GetOrder validate.Var: %v", err)
		return nil, err
	}

	order, err := d.pg.GetOrder(orderId)
	if err != nil {
		if _, ok := err.(*models.EntityError); ok {
			return nil, err
		}
		d.log.Error(err)
		return nil, models.ErrGeneralDBFail
	}

	return order, nil
}

// GetCustomerOrders gets orders for provided customer id. Returns EntityError if
// order or customer was not found and ErrGeneralDBFail if db returned db-specific error
func (d *dvdstoreUC) GetCustomerOrders(customerId int) ([]*models.Order, error) {
	if err := validateVar(customerId, "customerId"); err != nil {
		d.log.Debugf("GetCustomerOrders validate.Var: %v", err)
		return nil, err
	}

	// Check if customer exists
	var entErr *models.EntityError
	_, err := d.GetCustomer(customerId)
	if err != nil {
		if errors.As(err, &entErr) {
			return nil, err
		}
		d.log.Error(err)
		return nil, models.ErrGeneralDBFail
	}

	// Get orders
	orders, err := d.pg.GetCustomerOrders(customerId)
	if err != nil {
		if errors.As(err, &entErr) {
			return nil, err
		}
		d.log.Error(err)
		return nil, models.ErrGeneralDBFail
	}

	return orders, nil
}

// AddOrder creates order for customerId with provided products. Returns order and errors:
// EntityError if product/customer was not found or product is out of inventory
// and ErrGeneralDBFail if db returned db-specific error
func (d *dvdstoreUC) AddOrder(customerId int, products []*models.Product) (*models.Order, error) {
	// Validate inputs
	if err := validateVar(customerId, "customerId"); err != nil {
		d.log.Debugf("AddOrder validate.Var: %v", err)
		return nil, err
	}
	if len(products) == 0 {
		return nil, errors.New("products must not be empty")
	}
	for _, p := range products {
		if err := d.validate.StructPartial(p, "Id", "Quantity"); err != nil {
			d.log.Debugf("AddOrder validate.StructPartial: %v", err)
			return nil, models.ErrFieldsNotValid("id", "quantity")
		}
	}

	// Check if customer exists
	_, err := d.GetCustomer(customerId)
	var entErr *models.EntityError
	if err != nil {
		if errors.As(err, &entErr) {
			return nil, err
		}
		d.log.Error(err)
		return nil, models.ErrGeneralDBFail
	}

	// Add order
	order, err := d.pg.AddOrder(customerId, products)
	if err != nil {
		if errors.As(err, &entErr) {
			return nil, err
		}
		d.log.Error(err)
		return nil, models.ErrGeneralDBFail
	}

	return order, nil
}

// DeleteOrder deletes order by given order id and ErrGeneralDBFail if db returned db-specific error
func (d *dvdstoreUC) DeleteOrder(orderId int) error {
	if err := validateVar(orderId, "orderId"); err != nil {
		d.log.Debugf("DeleteOrder validate.Var: %v", err)
		return err
	}

	err := d.pg.DeleteOrder(orderId)
	if err != nil {
		d.log.Error(err)
		return models.ErrGeneralDBFail
	}
	return nil
}

// validateVar is a helper function that checks if var is > 0 and returns
// ValidationError if not ok
func validateVar(variable int, varName string) error {
	if variable > 0 {
		return nil
	}

	return &models.ValidationError{
		Message: fmt.Sprintf("%v must be > 0", varName),
	}
}
