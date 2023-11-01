package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"sync"
)

type Validator interface {
	GetValidate() *validator.Validate
	Validate(i interface{}) error
	AddTranslation(tag string, translation string) error
}

// customValidator is a struct for validator.
type customValidator struct {
	validate *validator.Validate
}

var validatorOnce sync.Once
var myValidator Validator

func GetValidator() (Validator, error) {
	validatorOnce.Do(func() {
		validate := validator.New()
		if err := registerTranslation(validate); err != nil {
			return
		}
		myValidator = &customValidator{
			validate: validate,
		}
	})

	return myValidator, nil
}

// GetValidate is a function for get validator.
func (v *customValidator) GetValidate() *validator.Validate {
	return v.validate
}

// Validate is a method for validate.
func (v *customValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func (v *customValidator) AddTranslation(tag string, translation string) error {
	return v.validate.RegisterTranslation(tag, GetTranslator(), func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())

		return t
	})
}

// registerTranslation is a function for register translation.
func registerTranslation(validate *validator.Validate) (err error) {
	if err = en_translations.RegisterDefaultTranslations(validate, GetTranslator()); err != nil {
		return
	}

	if err = validate.RegisterTranslation("required", GetTranslator(), func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	}); err != nil {
		return
	}

	return
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
