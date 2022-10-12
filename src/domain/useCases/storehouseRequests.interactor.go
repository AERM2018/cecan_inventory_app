package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"fmt"

	"github.com/kataras/iris/v12"
)

type StorehouseRequestsInteractor struct {
	StorehouseRequestDataSource datasources.StorehouseRequestsDataSource
}

func (interactor StorehouseRequestsInteractor) GetStorehouseRequests() models.Responser {
	storehouseRequests, err := interactor.StorehouseRequestDataSource.GetStorehouseRequests()
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"requests": storehouseRequests,
		},
	}
}

func (interactor StorehouseRequestsInteractor) GetStorehouseRequestById(id string) models.Responser {
	storehouseRequests, err := interactor.StorehouseRequestDataSource.GetStorehouseRequestById(id)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"request": storehouseRequests,
		},
	}
}

func (interactor StorehouseRequestsInteractor) CreateStorehouseRequest(storehouseRequest models.StorehouseRequestDetailed) models.Responser {
	storehouseRequestNoUtilities := models.StorehouseRequest{
		UserId: storehouseRequest.UserId,
	}
	fmt.Println(storehouseRequestNoUtilities)
	storehouseRequestId, errCreating := interactor.StorehouseRequestDataSource.CreateStorehouseRequest(storehouseRequestNoUtilities, storehouseRequest.Utilities)
	if errCreating != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        errCreating,
		}
	}
	storehouseRequest, errGetting := interactor.StorehouseRequestDataSource.GetStorehouseRequestById(storehouseRequestId.String())
	if errGetting != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        errGetting,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"request": storehouseRequest,
		},
	}
}
