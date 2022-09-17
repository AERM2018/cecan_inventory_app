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

type PharmacyStocksController struct {
	PharmacyStocksDataSource datasources.PharmacyStocksDataSource
	MedicineDataSource       datasources.MedicinesDataSource
	PharmacyStocksInteractor usecases.PharmacyStockInteractor
}

func (controller *PharmacyStocksController) New() {
	controller.PharmacyStocksInteractor = usecases.PharmacyStockInteractor{
		PharmacyStocksDataSource: controller.PharmacyStocksDataSource,
		MedicinesDataSource:      controller.MedicineDataSource,
	}
}

func (controller PharmacyStocksController) GetPharmacyStocks(ctx iris.Context) {
	res := controller.PharmacyStocksInteractor.GetPharmacyStocks()
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller PharmacyStocksController) UpdatePharmacyStock(ctx iris.Context) {
	var pharmacyStock models.PharmacyStockToUpdate
	pharmacyStockId := ctx.Params().GetString("id")
	id, _ := uuid.Parse(pharmacyStockId)
	bodyreader.ReadBodyAsJson(ctx, &pharmacyStock, true)
	res := controller.PharmacyStocksInteractor.UpdatePharmacyStock(id, pharmacyStock)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller PharmacyStocksController) DeletePharmacyStock(ctx iris.Context) {
	pharmacyStockId := ctx.Params().GetString("id")
	pharmacyStockUuid, _ := uuid.Parse(pharmacyStockId)
	res := controller.PharmacyStocksInteractor.DeletePharmacyStock(pharmacyStockUuid, true)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
