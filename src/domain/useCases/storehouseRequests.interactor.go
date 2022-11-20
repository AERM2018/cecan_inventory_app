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

func (interactor StorehouseRequestsInteractor) GetStorehouseRequestById(id string, isPdf bool) models.Responser {
	storehouseRequests, err := interactor.StorehouseRequestDataSource.GetStorehouseRequestById(id)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	if isPdf {
		storehouseRequestDoc := models.StorehouseRequestDoc{StorehouseRequest: storehouseRequests}
		filePath, errInPdf := storehouseRequestDoc.CreateDoc()
		if errInPdf != nil {
			return models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    "Ocurrió un error al generar el documento digital, intentelo mas tarde.",
			}
		}
		return models.Responser{
			ExtraInfo: []map[string]interface{}{
				{"file": filePath},
			},
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
			StatusCode: iris.StatusBadRequest,
			Err:        errCreating,
			Message:    errCreating.Error(),
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

func (interacor StorehouseRequestsInteractor) UpdateStorehouseRequest(id string, storehouseRequest models.StorehouseRequestDetailed) models.Responser {
	requestInfo := models.StorehouseRequest{}
	// Get the previous state of the storehouse request
	oldStorehouseRequest, _ := interacor.StorehouseRequestDataSource.GetStorehouseRequestById(id)
	// Get which utilities will be added, removed or just updated
	utilitiesToAdd, utilitiesToRemove := storehouseRequest.FilterUtilitesFromRequest(oldStorehouseRequest.Utilities)
	utilities := [][]models.StorehouseUtilitiesStorehouseRequests{storehouseRequest.Utilities, utilitiesToAdd, utilitiesToRemove}
	_, errUpdatingUtilities := interacor.StorehouseRequestDataSource.UpdateStorehouseRequest(id, requestInfo, utilities...)
	if errUpdatingUtilities != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        errUpdatingUtilities,
		}
	}
	storehouseRequest, errGettingRequest := interacor.StorehouseRequestDataSource.GetStorehouseRequestById(id)
	if errGettingRequest != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        errUpdatingUtilities,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"request": storehouseRequest,
		},
	}
}

func (interactor StorehouseRequestsInteractor) SupplyStorehouseRequest(id string, utilities []models.StorehouseUtilitiesStorehouseRequests) models.Responser {
	err := interactor.StorehouseRequestDataSource.SupplyStorehouseRequest(id, utilities)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Err:        err,
			Message:    err.Error(),
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Message:    "Los elementos de almacen fueron suministrados correctamente.",
	}
}
func (interactor StorehouseRequestsInteractor) DeleteStorehouseRequest(id string) models.Responser {
	err := interactor.StorehouseRequestDataSource.DeleteStorehouseRequest(id)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusNoContent,
	}
}
