package endpoints

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
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

func CustomDateValidator(fl validator.FieldLevel) bool {
	dateString := fl.Field().String()

	// Разделяем строку на компоненты диапазона, если есть ":"
	dateComponents := strings.Split(dateString, ":")
	if len(dateComponents) > 2 {
		return false
	}

	// Проверяем каждую компоненту на соответствие формату даты
	for _, date := range dateComponents {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			return false
		}
	}

	// Если все компоненты прошли проверку, возвращаем true
	return true
}

func CustomDateRangeValidator(fl validator.FieldLevel) bool {
	dateString := fl.Field().String()

	// Разделяем строку на компоненты диапазона, если есть ":"
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
