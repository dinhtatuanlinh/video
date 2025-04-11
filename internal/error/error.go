package error

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"os"
)

type AppError struct {
	Message  string
	HTTPCode int
	GRPCCode codes.Code
	Internal error
}

func (e *AppError) Error() string {
	env := os.Getenv("ENV")
	if e.Internal != nil {
		if env == "local" {
			return fmt.Sprintf("%s | internal: %v", e.Message, e.Internal)
		} else if env == "prod" {
			return e.Message
		}
		return fmt.Sprintf("%s | internal: %v", e.Message, e.Internal)
	}
	return e.Message
}

// Factory functions
func NewAppError(message string, httpCode int, grpcCode codes.Code, internal error) *AppError {
	return &AppError{
		Message:  message,
		HTTPCode: httpCode,
		GRPCCode: grpcCode,
		Internal: internal,
	}
}

func Wrap(err error, message string, httpCode int, grpcCode codes.Code) *AppError {
	return &AppError{
		Message:  message,
		HTTPCode: httpCode,
		GRPCCode: grpcCode,
		Internal: err,
	}
}
