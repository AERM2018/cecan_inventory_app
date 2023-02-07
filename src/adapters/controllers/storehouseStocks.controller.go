package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type StorehouseStocksController struct {
	StorehouseStocksDataSource    datasources.StorehouseStocksDataSource
	StorehouseUtilitiesDataSource datasources.StorehouseUtilitiesDataSource
	Interactor                    usecases.StorehouseStocksInteractor
}

func (controller *StorehouseStocksController) New() {
	controller.Interactor = usecases.StorehouseStocksInteractor{
		StorehouseStocksDataSource:    controller.StorehouseStocksDataSource,
		StorehouseUtilitiesDataSource: controller.StorehouseUtilitiesDataSource,
	}
}

func (controller StorehouseStocksController) GetStorehouseInventory(ctx iris.Context) {
	utilityKey := ctx.URLParam("utility_key")
	utilityName := ctx.URLParam("utility_name")
	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)
	includeDeleted, errBoolParse := ctx.URLParamBool("include_deleted")
	if errBoolParse != nil {
		includeDeleted = false
	}
	showLessQty, errBoolParse2 := ctx.URLParamBool("show_less_qty")
	if errBoolParse2 != nil {
		includeDeleted = false
	}
	// Build filters struct
	utilityFilters := models.StorehouseUtilitiesFilters{
		UtilityKey:     utilityKey,
		UtilityName:    utilityName,
		Page:           page,
		Limit:          limit,
		IncludeDeleted: includeDeleted,
		ShowLessQty:    showLessQty,
	}
	res := controller.Interactor.GetStorehouseInventory(utilityFilters)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseStocksController) UpdateStorehouseStock(ctx iris.Context) {
	var stock models.StorehouseStock
	storehouseStockId := ctx.Params().GetStringDefault("id", "")
	bodyreader.ReadBodyAsJson(ctx, &stock, true)
	res := controller.Interactor.UpdateStorehouseStock(storehouseStockId, stock)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller StorehouseStocksController) DeleteStorehouseStock(ctx iris.Context) {
	storehouseStockId := ctx.Params().GetStringDefault("id", "")
	res := controller.Interactor.DeleteStorehouseStock(storehouseStockId)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
