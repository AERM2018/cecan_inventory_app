package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type MedicinesController struct {
	MedicinesDataSource      datasources.MedicinesDataSource
	PharmacyStocksDataSource datasources.PharmacyStocksDataSource
	Interactor               usecases.MedicinesInteractor
}

func (controller *MedicinesController) New() {
	controller.Interactor = usecases.MedicinesInteractor{
		MedicinesDataSource:      controller.MedicinesDataSource,
		PharmacyStocksDataSource: controller.PharmacyStocksDataSource,
	}
}

func (controller MedicinesController) InsertMedicineIntoCatalog(ctx iris.Context) {
	medicine := models.Medicine{}
	ctx.ReadBody(&medicine)
	res := controller.Interactor.InsertMedicineIntoCatalog(medicine)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller MedicinesController) InsertPharmacyStockOfMedicine(ctx iris.Context) {
	var pharmacyStock models.PharmacyStock
	bodyreader.ReadBodyAsJson(ctx, &pharmacyStock, true)
	medicineKey := ctx.Params().GetString("key")
	pharmacyStock.MedicineKey = medicineKey
	res := controller.Interactor.InsertStockOfMedicine(pharmacyStock)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller MedicinesController) GetMedicinesCatalog(ctx iris.Context) {
	includeDeleted, _ := ctx.URLParamBool("include_deleted")
	medicineKey := ctx.URLParamDefault("key", "")
	medicineName := ctx.URLParamDefault("name", "")
	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)
	medicinesFilters := models.MedicinesFilters{
		MedicineKey:  medicineKey,
		MedicineName: medicineName,
		Limit:        limit,
		Page:         page,
	}
	res := controller.Interactor.GetMedicinesCatalog(medicinesFilters, includeDeleted)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller MedicinesController) UpdateMedicine(ctx iris.Context) {
	var medicine models.Medicine
	medicineKey := ctx.Params().GetString("key")
	bodyreader.ReadBodyAsJson(ctx, &medicine, true)
	res := controller.Interactor.UpdateMedicine(medicineKey, medicine)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller MedicinesController) DeleteMedicine(ctx iris.Context) {
	medicineKey := ctx.Params().GetString("key")
	res := controller.Interactor.DeleteMedicine(medicineKey)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller MedicinesController) ReactivateMedicine(ctx iris.Context) {
	medicineKey := ctx.Params().GetString("key")
	res := controller.Interactor.ReactivateMedicine(medicineKey)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
