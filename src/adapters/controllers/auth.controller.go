package controllers

import (
	"cecan_inventory/src/adapters/helpers"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"
	"cecan_inventory/src/infrastructure/storage/models"
	"errors"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type AuthController struct {
	UserDataSource datasources.UserDataSource
}
type AccessCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (controller AuthController) Login(ctx iris.Context) {
	credentials := AccessCredentials{}
	ctx.ReadBody(&credentials)
	user, err := controller.UserDataSource.GetUserByEmail(credentials.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.PrepareAndSendMessageResponse(ctx, iris.StatusNotFound, nil, "Invalid credentials.")
			return
		}
	}
	if isCorrectPassword := user.CheckPassword(credentials.Password); !isCorrectPassword {
		helpers.PrepareAndSendMessageResponse(ctx, iris.StatusNotFound, nil, "Invalid credentials.-")
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, iris.StatusOK, iris.Map{"user": user.ToJSON()})
}

func (controller AuthController) SignUp(ctx iris.Context) {
	var newUser models.User
	ctx.ReadBody(&newUser)
	newUserRecord, err := controller.UserDataSource.CreateUser(newUser)
	if err != nil {
		helpers.PrepareAndSendMessageResponse(ctx, iris.StatusNotFound, nil, err.Error())
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, iris.StatusOK, iris.Map{"user": newUserRecord})
}
