package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type StorehouseRequestsController struct {
	StorehouseRequestsDataSource datasources.StorehouseRequestsDataSource
	interactor                   usecases.StorehouseRequestsInteractor
}

func (controller *StorehouseRequestsController) New() {
	controller.interactor = usecases.StorehouseRequestsInteractor{
		StorehouseRequestDataSource: controller.StorehouseRequestsDataSource,
	}
}

func (controller StorehouseRequestsController) CreateStorehouseRequest(ctx iris.Context) {
	var storehouseRequest models.StorehouseRequestDetailed
	bodyreader.ReadBodyAsJson(ctx, &storehouseRequest, true)
	userId := ctx.Values().GetString("userId")
	storehouseRequest.UserId = userId
	res := controller.interactor.CreateStorehouseRequest(storehouseRequest)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
