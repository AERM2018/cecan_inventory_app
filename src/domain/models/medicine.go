package models

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

type Medicine struct {
	Key       string    `gorm:"primaryKey" json:"key" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at"`
}

func (medicine Medicine) Validate() (Responser, error) {
	customVal := validator.New()
	err := customVal.Struct(medicine)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorsMsgs := make(map[string]string)
		for _, valErr := range validationErrors {
			errorsMsgs[valErr.StructField()] = valErr.Error()
		}
		res := Responser{
			StatusCode: iris.StatusBadRequest,
			Data:       iris.Map{"errors": errorsMsgs},
		}
		return res, errors.New("Medicine struct validation error")
	}
	return Responser{}, nil
}
