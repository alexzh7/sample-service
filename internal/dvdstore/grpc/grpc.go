package grpc

import (
	"context"

	"github.com/alexzh7/sample-service/internal/dvdstore"
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
	d.log.Infof("Recieved GetCustomers call with limit %v", req.GetLimit())

	//try req.Limit

	// Get customers
	customers, err := d.uc.GetCustomers(int(req.GetLimit()))
	if err != nil {
		return nil, err //TODO: do err
	}

	// Customer model to grpc model
	customersProto := make([]*proto.Customer, 0)
	for _, c := range customers {
		cst := proto.Customer{
			Id:        int64(c.Id),
			FirstName: c.FirstName,
			LastName:  c.LastName,
			Age:       int64(c.Age),
		}
		customersProto = append(customersProto, &cst)
	}

	return &proto.GetCustomersRes{CustomerList: customersProto}, nil
}

// GetCustomer returns Customer by provided id
func (d *dvdstoreService) GetCustomer(ctx context.Context, req *proto.GetCustomerReq) (*proto.GetCustomerRes, error) {
	return nil, nil
}

// AddCustomer adds passed Customer and returns his id
func (d *dvdstoreService) AddCustomer(ctx context.Context, req *proto.AddCustomerReq) (*proto.AddCustomerRes, error) {
	return nil, nil
}

// DeleteCustomer deletes Customer by provided id
func (d *dvdstoreService) DeleteCustomer(ctx context.Context, req *proto.DeleteCustomerReq) (*proto.DeleteCustomerRes, error) {
	return nil, nil
}

// GetProducts returns list of all Products limited by provided limit
func (d *dvdstoreService) GetProducts(ctx context.Context, req *proto.GetProductsReq) (*proto.GetProductsRes, error) {
	return nil, nil
}

// GetProduct returns Product by provided id
func (d *dvdstoreService) GetProduct(ctx context.Context, req *proto.GetProductReq) (*proto.GetProductRes, error) {
	return nil, nil
}

// AddProduct adds passed Product and returns his id
func (d *dvdstoreService) AddProduct(ctx context.Context, req *proto.AddProductReq) (*proto.AddProductRes, error) {
	return nil, nil
}

// DeleteProduct deletes Product by provided id
func (d *dvdstoreService) DeleteProduct(ctx context.Context, req *proto.DeleteProductReq) (*proto.DeleteProductRes, error) {
	return nil, nil
}

// GetOrder gets order by provided id
func (d *dvdstoreService) GetOrder(ctx context.Context, req *proto.GetOrderReq) (*proto.GetOrderRes, error) {
	return nil, nil
}

// GetCustomerOrders returns customer orders by provided customer id
func (d *dvdstoreService) GetCustomerOrders(ctx context.Context, req *proto.GetCustomerOrdersReq) (*proto.GetCustomerOrdersRes, error) {
	return nil, nil
}

// AddOrder adds order for passed customer id with provided products and returns created order id
func (d *dvdstoreService) AddOrder(ctx context.Context, req *proto.AddOrderReq) (*proto.AddOrderRes, error) {
	// AddOrder returns only orderID
	return nil, nil
}

// DeleteOrder deletes order with provided order id
func (d *dvdstoreService) DeleteOrder(ctx context.Context, req *proto.DeleteOrderReq) (*proto.DeleteOrderRes, error) {
	return nil, nil
}
