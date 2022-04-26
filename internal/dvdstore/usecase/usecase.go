package usecase

import "github.com/alexzh7/sample-service/internal/dvdstore"

// dvdstoreUC is a dvd store use case
type dvdstoreUC struct {
	pg dvdstore.PostgresRepo
}

// Validation, logger, json errorrs, etc
