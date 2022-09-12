package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	structvalidator "cecan_inventory/infrastructure/external/structValidator"

	"strings"

	"github.com/kataras/iris/v12"
)

type AuthController struct {
	UserDataSource datasources.UserDataSource
	Interator      usecases.AuthInteractor
}

func (controller *AuthController) New(userDatasource datasources.UserDataSource) {
	controller.UserDataSource = userDatasource
	controller.Interator = usecases.AuthInteractor{UserDataSource: userDatasource}
}

func (controller AuthController) Login(ctx iris.Context) {
	credentials := models.AccessCredentials{}
	ctx.ReadBody(&credentials)
	res := controller.Interator.LoginUser(credentials)
	if res.StatusCode != iris.StatusOK {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller AuthController) SignUp(ctx iris.Context) {
	var newUser models.User
	ctx.ReadBody(&newUser)
	valRes, err := structvalidator.ValidateStructFomRequest(newUser)
	if err != nil {
		helpers.PrepareAndSendDataResponse(ctx, valRes)
		return
	}
	res := controller.Interator.SignUpUser(newUser)
	if res.Err != nil {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller AuthController) RenewToken(ctx iris.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	res := controller.Interator.RenewToken(token)
	if res.Err != nil {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
