package grpc

import "github.com/alexzh7/sample-service/internal/dvdstore"

// dvdstoreService is a grpc service for dvd store
// implements grpc server interface
type dvdstoreService struct {
	uc dvdstore.Usecase
}

// log, validate
