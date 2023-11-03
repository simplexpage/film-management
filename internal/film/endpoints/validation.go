package endpoints

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

// Regexp for validating date
var validDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// DateValidator Custom validator function for date
func DateValidator(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	if len(date) == 0 {
		return false
	}

	if !validDate.MatchString(date) {
		return false
	}

	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
