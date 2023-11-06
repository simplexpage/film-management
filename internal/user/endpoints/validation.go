package endpoints

import (
	"github.com/go-playground/validator/v10"
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
