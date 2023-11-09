package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
)

// ValidationError implements the Error interface.
type ValidationError struct {
	Field string
	Err   error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("field %s: err %v", e.Field, e.Err)
}

// NotFoundError implements the Error interface.
type NotFoundError struct {
	Err error
}

func (e NotFoundError) Error() string {
	return e.Err.Error()
}

// AuthError implements the Error interface.
type AuthError struct {
	Err error
}

func (e AuthError) Error() string {
	return e.Err.Error()
}

// CorsError implements the Error interface.
type CorsError struct {
	Err error
}

func (e CorsError) Error() string {
	return e.Err.Error()
}

// PermissionError implements the Error interface.
type PermissionError struct {
	Err error
}

func (e PermissionError) Error() string {
	return e.Err.Error()
}
