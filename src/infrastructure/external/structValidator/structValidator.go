package structvalidator

import (
	"cecan_inventory/domain/models"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

func ValidateStructFomRequest(structToValidate interface{}) (models.Responser, error) {
	customVal := validator.New()
	customVal.RegisterValidation("gttoday", isGraterThanToday)
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

func isGraterThanToday(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02 15:04:05 -0700 UTC", fmt.Sprintf("%v", fl.Field().Interface()))
	if err != nil {
		panic(err)
	}
	return date.Unix() > time.Now().Unix()
}
