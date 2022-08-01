package customvalidator

import "github.com/go-playground/validator/v10"

var customValidator = validator.New()

func GetValidatorInstance() *validator.Validate {
	if customValidator == nil {
		customValidator = validator.New()
	}
	return customValidator
}
