package validation

import (
	"fmt"
)

// CustomError implements the Error interface.
type CustomError struct {
	Field string
	Err   error
}

func (e CustomError) Error() string {
	return fmt.Sprintf("field %s: err %v", e.Field, e.Err)
}

// NotFoundError implements the Error interface.
type NotFoundError struct {
	Err error
}

func (e NotFoundError) Error() string {
	return e.Err.Error()
}
