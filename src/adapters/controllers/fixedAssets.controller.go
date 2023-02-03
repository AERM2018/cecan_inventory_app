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
	DepartmentsDataSource            datasources.DepartmentDataSource
	UserDataSource                   datasources.UserDataSource
	FixedAssetsInteractor            usecases.FixedAssetsInteractor
}

func (controller *FixedAssetsController) New() {
	controller.FixedAssetsInteractor = usecases.FixedAssetsInteractor{
		FixedAssetsDataSource:            controller.FixedAssetsDataSource,
		FixedAssetDescriptionsDataSource: controller.FixedAssetDescriptionsDataSource,
		DepartmentsDataSource:            controller.DepartmentsDataSource,
		UserDataSource:                   controller.UserDataSource,
	}
}

func (controller FixedAssetsController) GetFixedAssets(ctx iris.Context) {
	// This is for pagination
	page := ctx.URLParamIntDefault("page",1)
	limit := ctx.URLParamIntDefault("limit",10)
	offset := ctx.URLParamIntDefault("offset",limit)
	departmentName := ctx.URLParamDefault("department_name", "")
	brand := ctx.URLParamDefault("brand", "")
	model := ctx.URLParamDefault("model", "")
	typeOfAsset := ctx.URLParamDefault("type", "")
	physicState := ctx.URLParamDefault("physic_state", "")
	description := ctx.URLParamDefault("description", "")
	isPdf, _ := ctx.URLParamBool("pdf")
	fromDate := ctx.URLParamDefault("from", "2000/01/01")
	toDate := ctx.URLParamDefault("to", "9999/12/31")
	filters := models.FixedAssetFilters{
		DepartmentName: strings.ToUpper(departmentName),
		Brand:          strings.ToUpper(brand),
		Model:          strings.ToUpper(model),
		Type:           strings.ToUpper(typeOfAsset),
		PhysicState:    strings.ToUpper(physicState),
		Description:    strings.ToUpper(description),

	}
	datesDelimiter := []string{fmt.Sprintf("'%v'", fromDate), fmt.Sprintf("'%v'", toDate)}
	res := controller.FixedAssetsInteractor.GetFixedAssets(filters, datesDelimiter, isPdf, page, limit,offset)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	if isPdf {
		ctx.SendFile(fmt.Sprintf("%v", res.ExtraInfo[0]["file"]), "receta.pdf")
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

func (controller FixedAssetsController) UploadFixedAssetsFromFile(ctx iris.Context) {
	filePath, err := helpers.UploadFile(ctx, "seeds", "fixed_assets.csv","excel_file")
	if err != nil {
		helpers.PrepareAndSendMessageResponse(ctx, models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	res := controller.FixedAssetsInteractor.UploadFileDataToDb(filePath)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
