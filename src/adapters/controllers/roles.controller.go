package controllers

import (
	"cecan_inventory/adapters/helpers"
	usecases "cecan_inventory/domain/useCases"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type RolesController struct {
	RolesDataSource datasources.RolesDataSource
	Interactor      usecases.RolesInteractor
}

func (controller *RolesController) New() {
	controller.Interactor = usecases.RolesInteractor{RolesDataSource: controller.RolesDataSource}
}

func (controller RolesController) GetRoles(ctx iris.Context) {
	res := controller.Interactor.GetRoles()
	if res.StatusCode != iris.StatusOK {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
