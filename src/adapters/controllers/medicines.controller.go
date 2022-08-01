package controllers

import (
	"cecan_inventory/src/adapters/helpers"
	"cecan_inventory/src/domain/models"
	usecases "cecan_inventory/src/domain/useCases"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"

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
	valRes, err := medicine.Validate()
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
