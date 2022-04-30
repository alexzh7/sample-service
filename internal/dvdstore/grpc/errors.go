package grpc

import (
	"errors"

	"github.com/alexzh7/sample-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpcError returns valid errors for grpc
func grpcError(err error) error {
	return status.Errorf(getGrpcCode(err), err.Error())
}

// getGrpcCode assigns grpc error codes according to package errors
func getGrpcCode(err error) codes.Code {
	var validationErr *models.ValidationError
	var entityErr *models.EntityError
	switch {
	case errors.As(err, &validationErr):
		return codes.InvalidArgument
	case errors.As(err, &entityErr):
		return codes.NotFound
	case errors.Is(err, models.ErrGeneralDBFail):
		return codes.Internal
	}
	return codes.Internal
}
