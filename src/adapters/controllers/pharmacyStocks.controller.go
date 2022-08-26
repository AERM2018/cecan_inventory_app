package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	structvalidator "cecan_inventory/infrastructure/external/structValidator"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type PharmacyStocksController struct {
	PharmacyStocksDataSource datasources.PharmacyStocksDataSource
	PharmacyStocksInteractor usecases.PharmacyStockInteractor
}

func (controller *PharmacyStocksController) New(pharmacyStocksDataSource datasources.PharmacyStocksDataSource, medicineDataSource datasources.MedicinesDataSource) {
	controller.PharmacyStocksDataSource = pharmacyStocksDataSource
	controller.PharmacyStocksInteractor = usecases.PharmacyStockInteractor{PharmacyStocksDataSource: pharmacyStocksDataSource, MedicinesDataSource: medicineDataSource}
}

func (controller PharmacyStocksController) InsertStockOfMedicine(ctx iris.Context) {
	var stock models.PharmacyStock
	bodyreader.ReadBodyAsJson(ctx, &stock, true)
	valRes, err := structvalidator.ValidateStructFomRequest(stock)
	if err != nil {
		helpers.PrepareAndSendMessageResponse(ctx, valRes)
		return
	}
	res := controller.PharmacyStocksInteractor.InsertStockOfMedicine(stock)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
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
	bodyreader.ReadBodyAsJson(ctx, &pharmacyStock, true)
	id, _ := uuid.Parse(pharmacyStockId)
	res := controller.PharmacyStocksInteractor.UpdatePharmacyStock(id, pharmacyStock)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
