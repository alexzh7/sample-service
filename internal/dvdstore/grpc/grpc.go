package grpc

import (
	"context"

	"github.com/alexzh7/sample-service/internal/dvdstore"
	"github.com/alexzh7/sample-service/internal/models"
	"github.com/alexzh7/sample-service/proto"
	"go.uber.org/zap"
)

// dvdstoreService is a grpc service for dvd store. It implements grpc server interface
type dvdstoreService struct {
	uc  dvdstore.Usecase
	log *zap.SugaredLogger
	proto.UnimplementedDvdstoreServer
}

// NewDvdstoreService returns new dvd store service
func NewDvdstoreService(uc dvdstore.Usecase, log *zap.SugaredLogger) *dvdstoreService {
	return &dvdstoreService{uc: uc, log: log}
}

// GetCustomers returns list of all Customers limited by provided limit
func (d *dvdstoreService) GetCustomers(ctx context.Context, req *proto.GetCustomersReq) (*proto.GetCustomersRes, error) {
	limit := int(req.GetLimit())
	d.log.Infof("Received GetCustomers call with limit %v", limit)

	// Get customers
	customers, err := d.uc.GetCustomers(limit)
	if err != nil {
		return nil, grpcError(err)
	}

	// Form response
	customersProto := make([]*proto.Customer, 0)
	for _, c := range customers {
		customersProto = append(customersProto, c.ToProto())
	}

	return &proto.GetCustomersRes{CustomerList: customersProto}, nil
}

// GetCustomer returns Customer by provided id
func (d *dvdstoreService) GetCustomer(ctx context.Context, req *proto.GetCustomerReq) (*proto.GetCustomerRes, error) {
	customerId := int(req.GetCustomerID())
	d.log.Infof("Received GetCustomer call with id %v", customerId)

	customer, err := d.uc.GetCustomer(customerId)
	if err != nil {
		return nil, grpcError(err)
	}

	return &proto.GetCustomerRes{Customer: customer.ToProto()}, nil
}

// AddCustomer adds passed Customer and returns his id
func (d *dvdstoreService) AddCustomer(ctx context.Context, req *proto.AddCustomerReq) (*proto.AddCustomerRes, error) {
	d.log.Info("Received AddCustomer call")

	// Map proto.Customer to models.Customer
	cReq := req.GetCustomer()
	customer := &models.Customer{
		Id:        int(cReq.GetId()),
		FirstName: cReq.GetFirstName(),
		LastName:  cReq.GetLastName(),
		Age:       int(cReq.GetAge()),
	}

	id, err := d.uc.AddCustomer(customer)
	if err != nil {
		return nil, grpcError(err)
	}

	return &proto.AddCustomerRes{CustomerID: int64(id)}, nil
}

// DeleteCustomer deletes Customer by provided id
func (d *dvdstoreService) DeleteCustomer(ctx context.Context, req *proto.DeleteCustomerReq) (*proto.DeleteCustomerRes, error) {
	customerId := int(req.GetCustomerID())
	d.log.Infof("Received DeleteCustomer call with id %v", customerId)

	if err := d.uc.DeleteCustomer(customerId); err != nil {
		return nil, grpcError(err)
	}

	return &proto.DeleteCustomerRes{}, nil
}

// GetProducts returns list of all Products limited by provided limit
func (d *dvdstoreService) GetProducts(ctx context.Context, req *proto.GetProductsReq) (*proto.GetProductsRes, error) {
	limit := int(req.GetLimit())
	d.log.Infof("Received GetProducts call with limit %v", limit)

	// Get products
	products, err := d.uc.GetProducts(limit)
	if err != nil {
		return nil, grpcError(err)
	}

	// Form response
	productsRes := make([]*proto.Product, 0)
	for _, p := range products {
		productsRes = append(productsRes, p.ToProto())
	}

	return &proto.GetProductsRes{ProductList: productsRes}, nil
}

// GetProduct returns Product by provided id
func (d *dvdstoreService) GetProduct(ctx context.Context, req *proto.GetProductReq) (*proto.GetProductRes, error) {
	productId := int(req.GetProductID())
	d.log.Infof("Received GetProduct call with id %v", productId)

	product, err := d.uc.GetProduct(productId)
	if err != nil {
		return nil, grpcError(err)
	}

	return &proto.GetProductRes{Product: product.ToProto()}, nil
}

// AddProduct adds passed Product and returns his id
func (d *dvdstoreService) AddProduct(ctx context.Context, req *proto.AddProductReq) (*proto.AddProductRes, error) {
	d.log.Info("Received AddProduct call")

	product := models.ProductFromProto(req.GetProduct())
	id, err := d.uc.AddProduct(product)
	if err != nil {
		return nil, grpcError(err)
	}

	return &proto.AddProductRes{ProductID: int64(id)}, nil
}

// DeleteProduct deletes Product by provided id
func (d *dvdstoreService) DeleteProduct(ctx context.Context, req *proto.DeleteProductReq) (*proto.DeleteProductRes, error) {
	productId := int(req.GetProductID())
	d.log.Infof("Received DeleteProduct call with id %v", productId)

	if err := d.uc.DeleteProduct(productId); err != nil {
		return nil, grpcError(err)
	}

	return &proto.DeleteProductRes{}, nil
}

// GetOrder gets order by provided id
func (d *dvdstoreService) GetOrder(ctx context.Context, req *proto.GetOrderReq) (*proto.GetOrderRes, error) {
	orderId := int(req.GetOrderID())
	d.log.Infof("Received GetOrder call with id %v", orderId)

	order, err := d.uc.GetOrder(orderId)
	if err != nil {
		return nil, grpcError(err)
	}

	return &proto.GetOrderRes{Order: order.ToProto()}, nil
}

// GetCustomerOrders returns customer orders by provided customer id
func (d *dvdstoreService) GetCustomerOrders(ctx context.Context, req *proto.GetCustomerOrdersReq) (*proto.GetCustomerOrdersRes, error) {
	customerId := int(req.GetCustomerID())
	d.log.Infof("Received GetCustomerOrders call with id %v", customerId)

	// Get orders
	orders, err := d.uc.GetCustomerOrders(customerId)
	if err != nil {
		return nil, grpcError(err)
	}

	// Form response
	protoOrders := make([]*proto.Order, 0)
	for _, o := range orders {
		protoOrders = append(protoOrders, o.ToProto())
	}

	return &proto.GetCustomerOrdersRes{OrderList: protoOrders}, nil
}

// AddOrder adds order for passed customer id with provided products and returns created order id
func (d *dvdstoreService) AddOrder(ctx context.Context, req *proto.AddOrderReq) (*proto.AddOrderRes, error) {
	customerId := int(req.GetCustomerID())
	d.log.Infof("Received AddOrder call for customer id %v", customerId)

	// Form request
	products := make([]*models.Product, 0)
	for _, p := range req.GetProductList() {
		products = append(products, models.ProductFromProto(p))
	}

	order, err := d.uc.AddOrder(customerId, products)
	if err != nil {
		return nil, grpcError(err)
	}

	return &proto.AddOrderRes{OrderID: int64(order.Id)}, nil
}

// DeleteOrder deletes order with provided order id
func (d *dvdstoreService) DeleteOrder(ctx context.Context, req *proto.DeleteOrderReq) (*proto.DeleteOrderRes, error) {
	orderId := int(req.GetOrderID())
	d.log.Infof("Received DeleteOrder call with id %v", orderId)

	if err := d.uc.DeleteOrder(orderId); err != nil {
		return nil, grpcError(err)
	}

	return &proto.DeleteOrderRes{}, nil
}
