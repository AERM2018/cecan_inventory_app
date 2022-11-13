package controllers

import (
	"cecan_inventory/adapters/helpers"
	usecases "cecan_inventory/domain/useCases"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type FixedAssetDescriptionsController struct {
	FixedAssetDescriptionsDataSource datasources.FixedAssetDescriptionDataSource
	FixedAssetDescriptionsInteractor usecases.FixedAssetDescriptionsInteractor
}

func (controller *FixedAssetDescriptionsController) New() {
	controller.FixedAssetDescriptionsInteractor = usecases.FixedAssetDescriptionsInteractor{
		FixedAssetDescriptionsDataSource: controller.FixedAssetDescriptionsDataSource,
	}
}

func (controller FixedAssetDescriptionsController) GetFixedAssetDescriptions(ctx iris.Context) {
	query := ctx.URLParamDefault("q", "")
	res := controller.FixedAssetDescriptionsInteractor.GetFixedAssetDescriptions(query)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
