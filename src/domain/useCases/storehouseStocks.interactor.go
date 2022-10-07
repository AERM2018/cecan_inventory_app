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

func (interactor StorehouseStocksInteractor) GetStorehouseInventory() models.Responser {
	storehouseInventory := make([]models.StorehouseUtilityStocksDetailed, 0)
	// Get utilities from storehouse catalog
	storehouseUtilities := make([]models.StorehouseUtilityDetailed, 0)
	storehouseUtilities, _ = interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilities(true)
	// Look for each storehouse stokcs in inventory
	for _, storehouseUtility := range storehouseUtilities {
		utilityStocks, _ := interactor.StorehouseStocksDataSource.GetStorehouseStocksByUtiltyKey(storehouseUtility.Key)
		utilityStocksDetailed := models.StorehouseUtilityStocksDetailed{
			StorehouseUtility: storehouseUtility,
			StorehouseStocks:  utilityStocks,
		}
		utilityStocksDetailed.GetTotalQuantitiesLeft()
		storehouseInventory = append(
			storehouseInventory,
			utilityStocksDetailed,
		)
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"inventory": storehouseInventory,
		},
	}
}
