package models

import "time"

// Customer model
type Customer struct {
	Id        int      `json:"id,omitempty"`
	FirstName string   `json:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty"`
	Age       int      `json:"age,omitempty"`
	Orders    []*Order `json:"orders,omitempty"`
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

// Product model
type Product struct {
	Id       int     `json:"id,omitempty"`
	Title    string  `json:"title,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Quantity int     `json:"quantity,omitempty"`
}
