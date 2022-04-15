package models

//Customer model
type Customer struct {
	Id        int    `sql:"customerid" json:"id,omitempty"`
	FirstName string `sql:"firstname" json:"firstName,omitempty"`
	LastName  string `sql:"lastname" json:"lastName,omitempty"`
	Age       int    `sql:"age" json:"age,omitempty"`
}
