package structvalidator

import (
	"cecan_inventory/src/domain/models"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

func ValidateStructFomRequest(structToValidate any) (models.Responser, error) {
	customVal := validator.New()
	err := customVal.Struct(structToValidate)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorsMsgs := make(map[string]string)
		for _, valErr := range validationErrors {
			fmt.Println(valErr)
			errorsMsgs[valErr.StructField()] = valErr.Error()
		}
		res := models.Responser{
			StatusCode: iris.StatusBadRequest,
			Data:       iris.Map{"errors": errorsMsgs},
		}
		return res, errors.New("struct validation error")
	}
	return models.Responser{}, nil
}
