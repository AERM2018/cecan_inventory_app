package controllers

import (
	"cecan_inventory/adapters/helpers"
	usecases "cecan_inventory/domain/useCases"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type UsersController struct {
	UserDataSource datasources.UserDataSource
	Interator      usecases.UsersInteractor
}

func (controller *UsersController) New() {
	controller.Interator = usecases.UsersInteractor{UserDataSource: controller.UserDataSource}
}

func (controller UsersController) GetUsers(ctx iris.Context) {
	res := controller.Interator.GetUsers()
	if res.StatusCode != iris.StatusOK {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
