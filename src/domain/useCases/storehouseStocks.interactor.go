package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type StorehouseStocksInteractor struct {
	StorehouseStocksDataSource    datasources.StorehouseStocksDataSource
	StorehouseUtilitiesDataSource datasources.StorehouseUtilitiesDataSource
}

func (interactor StorehouseStocksInteractor) CreateStorehouseStock(stock models.StorehouseStock) models.Responser {
	stockId, err := interactor.StorehouseStocksDataSource.CreateStorehouseStock(stock)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"id": stockId,
		},
	}
}
