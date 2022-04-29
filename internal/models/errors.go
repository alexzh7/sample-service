package models

import (
	"errors"
	"fmt"
)

// ErrGeneralDBFail is used to hide db errors from client
var ErrGeneralDBFail = errors.New("unexpected database error")

// EntityError represents all errors that can be exposed to the user,
// for example "order/product/customer not found"
type EntityError struct {
	Entity  string
	Message string
}

func (e *EntityError) Error() string {
	return fmt.Sprintf("%v %v", e.Entity, e.Message)
}

// ErrNotFound composes "not found" errors for provided entities
func ErrNotFound(entity string, id int) *EntityError {
	return &EntityError{Entity: entity, Message: fmt.Sprintf("id %v not found", id)}
}

// ErrOutOfInventory composes "out of inventory" errors for provided entities
func ErrOutOfInventory(entity string, id int) *EntityError {
	return &EntityError{Entity: entity, Message: fmt.Sprintf("id %v is out of inventory", id)}
}

// ValidationError represents validation errors
type ValidationError struct {
	Message string
}

func (v *ValidationError) Error() string {
	return v.Message
}

// ErrFieldsNotValid composes validation error with provided fields
func ErrFieldsNotValid(fieldNames ...string) *ValidationError {
	return &ValidationError{Message: fmt.Sprintf("%q fields must be valid", fieldNames)}
}
