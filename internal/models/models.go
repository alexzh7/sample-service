package models

import (
	"time"

	"github.com/alexzh7/sample-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Customer model
type Customer struct {
	Id        int      `json:"id,omitempty"`
	FirstName string   `json:"firstName,omitempty" validate:"required,max=50"`
	LastName  string   `json:"lastName,omitempty" validate:"required,max=50"`
	Age       int      `json:"age,omitempty" validate:"required,gt=0,max=150"`
	Orders    []*Order `json:"orders,omitempty"`
}

// Map models.Customer to proto.Customer
func (c *Customer) ToProto() *proto.Customer {
	return &proto.Customer{
		Id:        int64(c.Id),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Age:       int64(c.Age),
	}
}

// Product model
type Product struct {
	Id       int     `json:"id,omitempty" validate:"required,gte=0,int"`
	Title    string  `json:"title,omitempty" validate:"required,max=50"`
	Price    float64 `json:"price,omitempty" validate:"required,gte=0,float"`
	Quantity int     `json:"quantity,omitempty" validate:"required,gte=0,int"`
}

// Map models.Product to proto.Product
func (p *Product) ToProto() *proto.Product {
	return &proto.Product{
		Id:       int64(p.Id),
		Title:    p.Title,
		Price:    p.Price,
		Quantity: int64(p.Quantity),
	}
}

// ProductFromProto maps proto.Product to models.Product
func ProductFromProto(product *proto.Product) *Product {
	return &Product{
		Id:       int(product.GetId()),
		Title:    product.GetTitle(),
		Price:    product.GetPrice(),
		Quantity: int(product.GetQuantity()),
	}
}

// Order model
type Order struct {
	Id          int        `json:"id,omitempty"`
	Date        time.Time  `json:"date,omitempty"`
	NetAmount   float64    `json:"netamount,omitempty"`
	Tax         float64    `json:"tax,omitempty"`
	TotalAmount float64    `json:"totalamount,omitempty"`
	Products    []*Product `json:"products,omitempty"`
}

// Map models.Order to proto.Order
func (o *Order) ToProto() *proto.Order {
	// Fill products if they exist
	products := make([]*proto.Product, 0)
	if len(o.Products) > 0 {
		for _, p := range o.Products {
			products = append(products, p.ToProto())
		}
	}

	return &proto.Order{
		Id:          int64(o.Id),
		Date:        timestamppb.New(o.Date),
		NetAmount:   o.NetAmount,
		Tax:         o.Tax,
		TotalAmount: o.TotalAmount,
		ProductList: products,
	}
}

// SortById implements sort.Interface and is used to sort slice of products by id
type SortById []*Product

func (a SortById) Len() int           { return len(a) }
func (a SortById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortById) Less(i, j int) bool { return a[i].Id < a[j].Id }
