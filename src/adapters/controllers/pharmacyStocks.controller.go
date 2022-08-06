package controllers

import (
	"cecan_inventory/src/adapters/helpers"
	"cecan_inventory/src/domain/models"
	usecases "cecan_inventory/src/domain/useCases"
	bodyreader "cecan_inventory/src/infrastructure/external/bodyReader"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"
	structvalidator "cecan_inventory/src/infrastructure/external/structValidator"

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
