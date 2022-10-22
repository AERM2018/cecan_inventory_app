package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/google/uuid"
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

func (controller StorehouseRequestsController) GetStorehouseRequests(ctx iris.Context) {
	res := controller.interactor.GetStorehouseRequests()
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseRequestsController) GetStorehouseRequestById(ctx iris.Context) {
	storehouseRequetId := ctx.Params().GetStringDefault("id", "")
	res := controller.interactor.GetStorehouseRequestById(storehouseRequetId)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseRequestsController) UpdateStorehouseRequest(ctx iris.Context) {
	var storehouseRequest models.StorehouseRequestDetailed
	storehouseRequetId := ctx.Params().GetStringDefault("id", "")
	uuidFromString, _ := uuid.Parse(storehouseRequetId)
	bodyreader.ReadBodyAsJson(ctx, &storehouseRequest, true)
	storehouseRequest.Id = uuidFromString
	res := controller.interactor.UpdateStorehouseRequest(storehouseRequetId, storehouseRequest)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseRequestsController) DeleteStorehouseRequest(ctx iris.Context) {
	storehouseRequetId := ctx.Params().GetStringDefault("id", "")
	res := controller.interactor.DeleteStorehouseRequest(storehouseRequetId)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
