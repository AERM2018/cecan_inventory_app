package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	structvalidator "cecan_inventory/infrastructure/external/structValidator"

	"github.com/kataras/iris/v12"
)

type MedicinesController struct {
	MedicinesDataSource datasources.MedicinesDataSource
	Interactor          usecases.MedicinesInteractor
}

func (controller *MedicinesController) New(medicinesDataSource datasources.MedicinesDataSource) {
	controller.MedicinesDataSource = medicinesDataSource
	controller.Interactor = usecases.MedicinesInteractor{MedicinesDataSource: medicinesDataSource}
}

func (controller MedicinesController) InsertMedicineIntoCatalog(ctx iris.Context) {
	medicine := models.Medicine{}
	ctx.ReadBody(&medicine)
	valRes, err := structvalidator.ValidateStructFomRequest(medicine)
	if err != nil {
		helpers.PrepareAndSendDataResponse(ctx, valRes)
		return
	}
	res := controller.Interactor.InsertMedicineIntoCatalog(medicine)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller MedicinesController) GetMedicinesCatalog(ctx iris.Context) {
	res := controller.Interactor.GetMedicinesCatalog()
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
