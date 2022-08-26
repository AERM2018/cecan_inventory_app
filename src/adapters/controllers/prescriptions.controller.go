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

type PrescriptionsController struct {
	PrescriptionsDataSource datasources.PrescriptionsDataSource
	PrescriptionsInteractor usecases.PrescriptionInteractor
}

func (controller *PrescriptionsController) New(prescriptionsDataSource datasources.PrescriptionsDataSource) {
	controller.PrescriptionsDataSource = prescriptionsDataSource
	controller.PrescriptionsInteractor = usecases.PrescriptionInteractor{PrescriptionsDataSource: prescriptionsDataSource}
}

func (controller PrescriptionsController) CreatePrescription(ctx iris.Context) {
	var prescriptionRequest models.PrescriptionDetialed
	bodyreader.ReadBodyAsJson(ctx, &prescriptionRequest, true)
	res := controller.PrescriptionsInteractor.CreatePrescription(prescriptionRequest)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller PrescriptionsController) GetPrescriptionById(ctx iris.Context) {
	idString := ctx.Params().GetString("id")
	id, _ := uuid.Parse(idString)
	res := controller.PrescriptionsInteractor.GetPrescriptionById(id)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
