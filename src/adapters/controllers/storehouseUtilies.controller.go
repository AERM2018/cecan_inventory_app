package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type StorehouseUtilitiesController struct {
	StorehouseUtilitiesDataSource datasources.StorehouseUtilitiesDataSource
	StorehouseStocksDataSource    datasources.StorehouseStocksDataSource
	Interactor                    usecases.StorehouseUtilitiesInteractor
}

func (controller *StorehouseUtilitiesController) New() {
	controller.Interactor = usecases.StorehouseUtilitiesInteractor{
		StorehouseUtilitiesDataSource: controller.StorehouseUtilitiesDataSource,
		StorehouseStocksDataSource:    controller.StorehouseStocksDataSource,
	}
}

func (controller StorehouseUtilitiesController) CreateStorehouseUtility(ctx iris.Context) {
	var storehouseUtility models.StorehouseUtility
	bodyreader.ReadBodyAsJson(ctx, &storehouseUtility, true)
	res := controller.Interactor.CreateStorehouseUtility(storehouseUtility)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseUtilitiesController) CreateStorehouseUtilityStock(ctx iris.Context) {
	var stock models.StorehouseStock
	storehouseUtilityKey := ctx.Params().GetStringDefault("key", "")
	bodyreader.ReadBodyAsJson(ctx, &stock, true)
	res := controller.Interactor.CreateStorehouseUtilityStock(storehouseUtilityKey, stock)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseUtilitiesController) GetStorehouseUtilities(ctx iris.Context) {
	includeDeleted, _ := ctx.URLParamBool("include_deleted")
	res := controller.Interactor.GetStorehouseUtilities(includeDeleted)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseUtilitiesController) GetStorehouseUtility(ctx iris.Context) {
	storehouseUtilityKey := ctx.Params().GetStringDefault("key", "")
	res := controller.Interactor.GetStorehouseUtilityByKey(storehouseUtilityKey)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseUtilitiesController) UpdateStorehouseUtility(ctx iris.Context) {
	var storehouseUtility models.StorehouseUtility
	storehouseUtilityKey := ctx.Params().GetStringDefault("key", "")
	bodyreader.ReadBodyAsJson(ctx, &storehouseUtility, true)
	res := controller.Interactor.UpdateStorehouseUtility(storehouseUtilityKey, storehouseUtility)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseUtilitiesController) ReactivateStorehouseUtility(ctx iris.Context) {
	storehouseUtilityKey := ctx.Params().GetStringDefault("key", "")
	res := controller.Interactor.ReactivateStorehouseUtility(storehouseUtilityKey)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseUtilitiesController) DeleteStorehouseUtility(ctx iris.Context) {
	storehouseUtilityKey := ctx.Params().GetStringDefault("key", "")
	res := controller.Interactor.DeleteStorehouseUtility(storehouseUtilityKey)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseUtilitiesController) GetStorehouseUtilityCategories(ctx iris.Context) {
	res := controller.Interactor.GetStorehouseUtilityCategories()
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseUtilitiesController) GetStorehouseUtilityPresentations(ctx iris.Context) {
	res := controller.Interactor.GetStorehouseUtilityPresentations()
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseUtilitiesController) GetStorehouseUtilityUnits(ctx iris.Context) {
	res := controller.Interactor.GetStorehouseUtilityUnits()
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}