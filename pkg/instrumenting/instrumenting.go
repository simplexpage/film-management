package instrumenting

import (
	"errors"
	customError "film-management/pkg/errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

// PrintErr returns a string representation of an error, suitable for instrumentation.
func PrintErr(err error) string {
	isError := false

	if err != nil {
		var (
			validationErrors validator.ValidationErrors
			validationError  customError.ValidationError
		)

		if errors.As(err, &validationErrors) || errors.As(err, &validationError) {
			isError = false
		} else {
			isError = true
		}
	}

	return fmt.Sprint(isError)
}
