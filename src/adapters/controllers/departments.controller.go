package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type DepartmentsController struct {
	DepartmentsDataSource datasources.DepartmentDataSource
	DepartmentsInteractor usecases.DepartmentsInteractor
}

func (controller *DepartmentsController) New() {
	controller.DepartmentsInteractor = usecases.DepartmentsInteractor{
		DepartmentsDataSource: controller.DepartmentsDataSource,
	}
}

func (controller DepartmentsController) GetDepartments(ctx iris.Context) {
	includeDeleted, _ := ctx.URLParamBool("include_deleted")
	limit := ctx.URLParamIntDefault("limit", 10)
	offset := ctx.URLParamIntDefault("offset", 0)
	res := controller.DepartmentsInteractor.GetDepartments(includeDeleted, limit, offset)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller DepartmentsController) CreateDepartment(ctx iris.Context) {
	var department models.Department
	bodyreader.ReadBodyAsJson(ctx, &department, true)
	res := controller.DepartmentsInteractor.CreateDepartment(department)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller DepartmentsController) UpdateDepartment(ctx iris.Context) {
	var department models.Department
	bodyreader.ReadBodyAsJson(ctx, &department, true)
	id := ctx.Params().GetStringDefault("id", "")
	res := controller.DepartmentsInteractor.UpdateDepartment(id, department)
	if res.StatusCode > 300 {
		helpers.PrepareAndSendMessageResponse(ctx, res)
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}
