package models

import "time"

// Customer model
//TODO: I really need sql tag?
type Customer struct {
	Id        int      `sql:"customerid" json:"id,omitempty"`
	FirstName string   `sql:"firstname" json:"firstName,omitempty"`
	LastName  string   `sql:"lastname" json:"lastName,omitempty"`
	Age       int      `sql:"age" json:"age,omitempty"`
	Orders    []*Order `json:"orders,omitempty"`
}

// Order model
type Order struct {
	Id          int        `sql:"orderid" json:"id,omitempty"`
	Date        time.Time  `sql:"orderdate" json:"date,omitempty"`
	NetAmount   float64    `sql:"netamount" json:"netamount,omitempty"`
	Tax         float64    `sql:"tax" json:"tax,omitempty"`
	TotalAmount float64    `sql:"totalamount" json:"totalamount,omitempty"`
	Products    []*Product `json:"products,omitempty"`
}

// Product model
type Product struct {
	Id       int     `sql:"prod_id" json:"id,omitempty"`
	Title    string  `sql:"title" json:"title,omitempty"`
	Price    float64 `sql:"price" json:"price,omitempty"`
	Quantity int     `sql:"quan_in_stock" json:"quantity,omitempty"`
}

// Order by orderid
// SELECT p.title, p.price,
// t.orderid, t.orderdate, t.netamount, t.tax, t.totalamount,
// t.prod_id, t.quantity
// FROM products p INNER JOIN

// (SELECT o.orderid, o.orderdate, o.netamount, o.tax, o.totalamount,
// ol.prod_id, ol.quantity
// FROM orders o INNER JOIN orderlines ol
// ON o.orderid = ol.orderid) t

// ON p.prod_id = t.prod_id

// WHERE orderid=1

//Orders by customer
// SELECT p.title, p.price,
// t.quantity, t.prod_id,  t.customerid, t.orderid,
// t.orderdate, t.netamount, t.tax, t.totalamount
// FROM products p INNER JOIN

// (SELECT ol.prod_id, ol.quantity, o.*
// FROM  orderlines ol INNER JOIN orders o
// ON ol.orderid = o.orderid) t

// ON p.prod_id = t.prod_id

// WHERE customerid=201

// Orders by customer joining customer table
// SELECT c.customerid,
// t.orderid, t.orderdate, t.netamount, t.tax, t.totalamount,
// t.prod_id, t.quantity
// FROM customers c INNER JOIN

//     (SELECT p.title, p.price, t1.*
//     FROM products p INNER JOIN

//         (SELECT ol.prod_id, ol.quantity, o.*
//         FROM orderlines ol INNER JOIN orders o
//         ON ol.orderid = o.orderid) t1

//     ON p.prod_id = t1.prod_id) t

// ON c.customerid = t.customerid

// where c.customerid = 201
