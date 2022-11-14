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

type FixedAssetsRequestsController struct {
	FixedAssetsRequestsDataSource     datasources.FixedAssetsRequetsDataSource
	FixedAssetsDataSource             datasources.FixedAssetsDataSource
	FixedAssetsDescriptionsDataSource datasources.FixedAssetDescriptionDataSource
	Interactor                        usecases.FixedAssetsRequestsInteractor
	FixedAssetsInteractor             usecases.FixedAssetsInteractor
}

func (controller *FixedAssetsRequestsController) New() {
	controller.Interactor = usecases.FixedAssetsRequestsInteractor{
		FixedAssetsRequestsDataSource: controller.FixedAssetsRequestsDataSource,
		FixedAssetsDataSource:         controller.FixedAssetsDataSource,
	}
	controller.FixedAssetsInteractor = usecases.FixedAssetsInteractor{
		FixedAssetsDataSource:            controller.FixedAssetsDataSource,
		FixedAssetDescriptionsDataSource: controller.FixedAssetsDescriptionsDataSource,
	}
}

func (controller FixedAssetsRequestsController) CreateFixedAssetsRequest(ctx iris.Context) {
	var fixedAssetsRequest models.FixedAssetsRequestDetailed
	bodyreader.ReadBodyAsJson(ctx, &fixedAssetsRequest, true)
	userId := ctx.Values().GetString("userId")
	departmentId := ctx.Params().GetStringDefault("departmentId", "")
	fixedAssetsRequest.UserId = userId
	departmentIdUuid, _ := uuid.Parse(departmentId)
	fixedAssetsRequest.DepartmentId = departmentIdUuid
	res := controller.Interactor.CreateFixedAssetsRequest(
		fixedAssetsRequest,
		controller.FixedAssetsInteractor.CreateFixedAsset,
		controller.FixedAssetsInteractor.DeleteFixedAsset,
	)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller FixedAssetsRequestsController) DeleteFixedAssetsRequest(ctx iris.Context) {
	fixedAssetsRequestId := ctx.Params().GetStringDefault("id", "")
	res := controller.Interactor.DeleteFixedAssetsRequest(fixedAssetsRequestId)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
