package repository

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexzh7/sample-service/internal/models"
)

var (
	// NotEqualErr is a helper function for assert.ObjectsAreEqual
	NotEqualErr = func(want, got interface{}) error {
		w, _ := json.MarshalIndent(want, "", "\t")
		g, _ := json.MarshalIndent(got, "", "\t")
		return fmt.Errorf("objects not equal, expected: %+v,\n got: %+v", string(w), string(g))
	}

	// Mock objects
	mockCustomer = &models.Customer{Id: 1, FirstName: "John", LastName: "Doe", Age: 40}
	mockProduct  = &models.Product{Id: 1, Title: "Interstellar", Price: 80.00, Quantity: 60}
	mockProducts = []*models.Product{
		{Id: 1, Title: "Interstellar", Price: 80.00, Quantity: 60},
		{Id: 2, Title: "John Wick", Price: 100.00, Quantity: 230},
		{Id: 3, Title: "Inception", Price: 120.00, Quantity: 400},
	}
)

// NewMock returns mock db connection
func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating mock db: %v", err)
	}
	return db, mock
}

// AnyTime is used for matching time
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
