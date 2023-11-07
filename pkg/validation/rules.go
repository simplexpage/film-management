package validation

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// UsernameValidator Custom validator function for username.
func UsernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) == 0 {
		return false
	}

	// Check if the first character is a letter
	if !unicode.IsLetter(rune(username[0])) {
		return false
	}

	// Check if all characters are alphanumeric
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}

	return true
}

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

// DateRangeValidator Custom validator function for date range
func DateRangeValidator(fl validator.FieldLevel) bool {
	dateString := fl.Field().String()

	dateComponents := strings.Split(dateString, ":")
	if len(dateComponents) > 2 {
		return false
	}

	for _, date := range dateComponents {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			return false
		}
	}

	return true
}

// DateRangeCorrectValidator Custom validator function for date range
func DateRangeCorrectValidator(fl validator.FieldLevel) bool {
	dateString := fl.Field().String()

	dateComponents := strings.Split(dateString, ":")
	if len(dateComponents) != 2 {
		// Если это не диапазон, сразу возвращаем true
		return true
	}

	firstDate, err1 := time.Parse("2006-01-02", dateComponents[0])
	secondDate, err2 := time.Parse("2006-01-02", dateComponents[1])

	if err1 != nil || err2 != nil || firstDate.After(secondDate) {
		return false
	}

	return true
}
