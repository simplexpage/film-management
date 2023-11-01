package validation

import (
	"fmt"
)

type CustomError struct {
	Field string
	Err   error
}

func (e CustomError) Error() string {
	return fmt.Sprintf("field %s: err %v", e.Field, e.Err)
}
