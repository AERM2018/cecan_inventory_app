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

func (interactor DepartmentsInteractor) GetDepartments(includeDeleted bool, limit int, offset int) models.Responser {
	departments, err := interactor.DepartmentsDataSource.GetDepartments(includeDeleted, limit, offset)
	var headers []map[string]interface{}
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	// Add header with the query params limit and offset set depending on the page
	if len(departments) == limit {
		headers = append(headers, map[string]interface{}{
			"name":  "next_page",
			"value": fmt.Sprintf("http://localhost:4000/api/v1/departments?limit=%v&offset=%v", limit, offset+limit),
		})
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"departments": departments,
		},
		Headers: headers,
		ExtraInfo: []map[string]interface{}{
			{"page": (offset / limit) + 1},
		},
	}
}

func (interactor DepartmentsInteractor) GetDepartmentById(id string) models.Responser {
	departments, err := interactor.DepartmentsDataSource.GetDepartmentById(id)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"departments": departments,
		},
	}
}

func (interactor DepartmentsInteractor) UpdateDepartment(id string, department models.Department) models.Responser {
	err := interactor.DepartmentsDataSource.UpdateDepartment(id, department)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	departmentFound, _ := interactor.DepartmentsDataSource.GetDepartmentById(id)
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"department": departmentFound,
		},
	}
}

func (interactor DepartmentsInteractor) AssingResponsibleToDepartment(id string, userId string) models.Responser {
	err := interactor.DepartmentsDataSource.AssingResponsibleToDepartment(id, userId)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	departmentFound, _ := interactor.DepartmentsDataSource.GetDepartmentById(id)
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"department": departmentFound,
		},
	}
}
