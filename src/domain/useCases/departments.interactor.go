package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"fmt"

	"github.com/kataras/iris/v12"
)

type DepartmentsInteractor struct {
	DepartmentsDataSource datasources.DepartmentDataSource
}

func (interactor DepartmentsInteractor) CreateDepartment(department models.Department) models.Responser {
	idDepartment, err := interactor.DepartmentsDataSource.CreateDepartment(department)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	departmentCreated, err := interactor.DepartmentsDataSource.GetDepartmentById(idDepartment.String())
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"department": departmentCreated,
		},
	}
}
