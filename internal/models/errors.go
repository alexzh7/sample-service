package models

import "errors"

var (
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrOrderNotFound         = errors.New("order not found")
	ErrProductNotFound       = errors.New("product not found")
	ErrProductOutOfInventory = errors.New("product out of inventory")

	// ErrGeneralDBFail is used to not expose db errors to client
	ErrGeneralDBFail = errors.New("unexpected database error")
)
