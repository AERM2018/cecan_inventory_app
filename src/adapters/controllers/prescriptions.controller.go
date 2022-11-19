package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"fmt"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type PrescriptionsController struct {
	PrescriptionsDataSource datasources.PrescriptionsDataSource
	PrescriptionsInteractor usecases.PrescriptionInteractor
}

func (controller *PrescriptionsController) New() {
	controller.PrescriptionsInteractor = usecases.PrescriptionInteractor{
		PrescriptionsDataSource: controller.PrescriptionsDataSource}
}

func (controller PrescriptionsController) CreatePrescription(ctx iris.Context) {
	var prescriptionRequest models.PrescriptionDetialed
	bodyreader.ReadBodyAsJson(ctx, &prescriptionRequest, true)
	prescriptionRequest.UserId = ctx.Values().GetString("userId")
	res := controller.PrescriptionsInteractor.CreatePrescription(prescriptionRequest)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller PrescriptionsController) GetPrescriptions(ctx iris.Context) {
	userId := ctx.URLParamDefault("user_id", "")
	res := controller.PrescriptionsInteractor.GetPrescriptions(userId)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller PrescriptionsController) GetPrescriptionById(ctx iris.Context) {
	var isPdf bool
	idString := ctx.Params().GetString("id")
	if ctx.URLParamDefault("pdf", "false") == "false" {
		isPdf = false
	} else {
		isPdf = true
	}
	id, _ := uuid.Parse(idString)
	res := controller.PrescriptionsInteractor.GetPrescriptionById(id, isPdf)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	if isPdf {
		ctx.SendFile(fmt.Sprintf("%v", res.ExtraInfo[0]["file"]), "receta.pdf")
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller PrescriptionsController) UpdatePrescription(ctx iris.Context) {
	var prescriptionRequest models.PrescriptionDetialed
	id := ctx.Params().GetString("id")
	bodyreader.ReadBodyAsJson(ctx, &prescriptionRequest, true)
	res := controller.PrescriptionsInteractor.UpdatePrescription(id, prescriptionRequest)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller PrescriptionsController) CompletePrescription(ctx iris.Context) {
	var prescription models.PrescriptionToComplete
	bodyreader.ReadBodyAsJson(ctx, &prescription, true)
	id := ctx.Params().GetString("id")
	res := controller.PrescriptionsInteractor.CompletePrescription(id, prescription)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller PrescriptionsController) DeletePrescription(ctx iris.Context) {
	id := ctx.Params().GetString("id")
	res := controller.PrescriptionsInteractor.DeletePrescription(id)
	if res.StatusCode >= 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
