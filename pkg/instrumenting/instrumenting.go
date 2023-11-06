package instrumenting

import (
	"errors"
	"film-management/pkg/validation"
	"fmt"
	"github.com/go-playground/validator/v10"
)

// PrintErr returns a string representation of an error, suitable for instrumentation.
func PrintErr(err error) string {
	isError := false

	if err != nil {
		var (
			validationErrors validator.ValidationErrors
			customError      validation.CustomError
		)

		if errors.As(err, &validationErrors) || errors.As(err, &customError) {
			isError = false
		} else {
			isError = true
		}
	}

	return fmt.Sprint(isError)
}
