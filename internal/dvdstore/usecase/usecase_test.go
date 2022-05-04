package usecase

import (
	"math"
	"testing"

	"github.com/alexzh7/sample-service/internal/models"
)

// TODO
// Implement dvdstore.PostgresRepo
// Check ALL validations

func TestAddProduct(t *testing.T) {

	val := models.NewValidation()

	product := &models.Product{
		Id:       1,
		Title:    "32",
		Price:    math.Inf(1),
		Quantity: 272237289213213,
	}

	err := val.Struct(product)

	t.Logf("\n\n ERR: %v \n\n", err)
}
