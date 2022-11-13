package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
)

type FixedAssetsController struct {
	FixedAssetsDataSource            datasources.FixedAssetsDataSource
	FixedAssetDescriptionsDataSource datasources.FixedAssetDescriptionDataSource
	FixedAssetsInteractor            usecases.FixedAssetsInteractor
}

func (controller *FixedAssetsController) New() {
	controller.FixedAssetsInteractor = usecases.FixedAssetsInteractor{
		FixedAssetsDataSource:            controller.FixedAssetsDataSource,
		FixedAssetDescriptionsDataSource: controller.FixedAssetDescriptionsDataSource,
	}
}

func (controller FixedAssetsController) GetFixedAssets(ctx iris.Context) {
	departmentName := ctx.URLParamDefault("department_name", "")
	filters := models.FixedAssetFilters{
		DepartmentName: strings.ToUpper(departmentName),
	}
	res := controller.FixedAssetsInteractor.GetFixedAssets(filters)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller FixedAssetsController) CreateFixedAsset(ctx iris.Context) {
	var fixedAsset models.FixedAsset
	bodyreader.ReadBodyAsJson(ctx, &fixedAsset, true)
	res := controller.FixedAssetsInteractor.CreateFixedAsset(fixedAsset)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller FixedAssetsController) UpdateFixedAsset(ctx iris.Context) {
	var fixedAsset models.FixedAsset
	bodyreader.ReadBodyAsJson(ctx, &fixedAsset, true)
	fmt.Println(fixedAsset)
	fixedAssetKey := ctx.Params().GetStringDefault("key", "")
	res := controller.FixedAssetsInteractor.UpdateFixedAsset(fixedAssetKey, fixedAsset)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller FixedAssetsController) DeleteFixedAsset(ctx iris.Context) {
	fixedAssetKey := ctx.Params().GetStringDefault("key", "")
	res := controller.FixedAssetsInteractor.DeleteFixedAsset(fixedAssetKey)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
