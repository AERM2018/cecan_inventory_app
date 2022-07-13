package controllers

import (
	"cecan_inventory/src/adapters/helpers"
	"cecan_inventory/src/domain/models"
	usecases "cecan_inventory/src/domain/useCases"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"

	iris "github.com/kataras/iris/v12"
)

type AuthController struct {
	UserDataSource datasources.UserDataSource
	Iterator usecases.UserInteractor
}

func (controller *AuthController) New(userDatasource datasources.UserDataSource){
	controller.UserDataSource = userDatasource
	controller.Iterator = usecases.UserInteractor{ UserDataSource:  userDatasource}
}

func (controller AuthController) Login(ctx iris.Context) {
	credentials := models.AccessCredentials{}
	ctx.ReadBody(&credentials)
	res := controller.Iterator.LoginUser(credentials)
	if(res.StatusCode != iris.StatusOK){
		helpers.PrepareAndSendMessageResponse(ctx,res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx,res)
}

func (controller AuthController) SignUp(ctx iris.Context) {
	var newUser models.User
	ctx.ReadBody(&newUser)
	res := controller.Iterator.SignUpUser(newUser)
	if res.Err != nil {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
