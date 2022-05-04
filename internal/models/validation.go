package models

import (
	"math"

	"github.com/go-playground/validator/v10"
)

// Validation embeds validator and has custom validations
type Validation struct {
	*validator.Validate
}

// NewValidation creates new validator and registers custom validations
func NewValidation() *Validation {
	v := validator.New()
	v.RegisterValidation("float", validFloat)
	v.RegisterValidation("int", validInt)

	return &Validation{Validate: v}
}

// validInt checks for valid integer
func validInt(field validator.FieldLevel) bool {
	fl := field.Field().Int()
	if fl > math.MaxInt32 {
		return false
	}
	return true
}

// validFloat checks valid float field
func validFloat(field validator.FieldLevel) bool {
	fl := field.Field().Float()
	if math.IsInf(fl, 0) || math.IsNaN(fl) {
		return false
	}

	return true
}
