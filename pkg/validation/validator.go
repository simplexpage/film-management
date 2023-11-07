package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"sync"
)

type Validator interface {
	Validate(i interface{}) error
}

// customValidator is a struct for validator.
type customValidator struct {
	validate *validator.Validate
}

// Validate is a method for validate.
func (v *customValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

var validatorOnce sync.Once
var myValidator Validator

func GetValidator() (Validator, error) {
	validatorOnce.Do(func() {
		if v, err := CreateNewValidator(); err != nil {
			myValidator = nil
		} else {
			myValidator = v
		}
	})

	return myValidator, nil
}

func CreateNewValidator() (Validator, error) {
	v := &customValidator{
		validate: validator.New(),
	}

	// Register translation
	if err := v.registerTranslation(); err != nil {
		return nil, err
	}

	// Register validation
	if err := v.registerValidation(); err != nil {
		return nil, err
	}

	return v, nil
}

// registerTranslation is a function for register translation.
func (v *customValidator) registerTranslation() (err error) {
	if err = en_translations.RegisterDefaultTranslations(v.validate, GetTranslator()); err != nil {
		return err
	}

	if err = v.addTranslation("required", "{0} is required"); err != nil {
		return err
	}

	if err = v.addTranslation("username", "{0} must be valid (alphanumeric starting with letter)"); err != nil {
		return err
	}

	if err = v.addTranslation("customDate", "{0} must be valid (YYYY-MM-DD)"); err != nil {
		return err
	}

	if err = v.addTranslation("customRangeDateCorrect", "{0} must be valid the first date must be less than the second date"); err != nil {
		return err
	}

	if err = v.addTranslation("customRangeDate", "{0} must be valid (YYYY-MM-DD or YYYY-MM-DD:YYYY-MM-DD)"); err != nil {
		return err
	}

	return
}

// RegisterValidation is a function for register validation.
func (v *customValidator) registerValidation() error {
	if err := v.validate.RegisterValidation("username", UsernameValidator); err != nil {
		return err
	}

	if err := v.validate.RegisterValidation("customDate", DateValidator); err != nil {
		return err
	}

	if err := v.validate.RegisterValidation("customRangeDate", DateRangeValidator); err != nil {
		return err
	}

	if err := v.validate.RegisterValidation("customRangeDateCorrect", DateRangeCorrectValidator); err != nil {
		return err
	}

	return nil
}

func (v *customValidator) addTranslation(tag string, translation string) error {
	return v.validate.RegisterTranslation(tag, GetTranslator(), func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())

		return t
	})
}

var uni *ut.UniversalTranslator
var translator ut.Translator
var translatorOnce sync.Once

// GetTranslator is a function for get translator.
func GetTranslator() ut.Translator {
	translatorOnce.Do(func() {
		enT := en.New()
		uni = ut.New(enT, enT)

		translator, _ = uni.GetTranslator("en")
	})

	return translator
}
