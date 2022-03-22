package grpcp

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorHandler returns a gRPC error if a gRPC problem or a regular error if not
// a gRPC error.
func ErrorHandler(err error) error {
	stat := status.Convert(err)
	if stat.Code() == codes.Unknown {
		return errors.New(stat.Message())
	}

	return err
}
