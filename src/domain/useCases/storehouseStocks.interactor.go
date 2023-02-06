package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"fmt"

	"github.com/kataras/iris/v12"
)

type StorehouseStocksInteractor struct {
	StorehouseStocksDataSource    datasources.StorehouseStocksDataSource
	StorehouseUtilitiesDataSource datasources.StorehouseUtilitiesDataSource
}

func (interactor StorehouseStocksInteractor) GetStorehouseInventory(filters models.StorehouseUtilitiesFilters) models.Responser {
	var (
		stocksError error
		totalPages  int
	)
	storehouseInventoryByStocks := make([]models.StorehouseUtilityStockDetailed, 0)
	storehouseInventoryByUtility := make([]models.StorehouseUtilitsDetailedNoStocks, 0)
	data := iris.Map{}
	// Get utilities from storehouse catalog
	// storehouseUtilities := make([]models.StorehouseUtilityDetailed, 0)
	// storehouseUtilities, _ = interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilities(true)
	// Look for each storehouse stokcs in inventory
	// for _, storehouseUtility := range storehouseUtilities {
	// 	utilityStocks, _ := interactor.StorehouseStocksDataSource.GetStorehouseStocksByUtiltyKey(storehouseUtility.Key)
	// 	utilityStocksDetailed := models.StorehouseUtilityStocksDetailed{
	// 		StorehouseUtility: storehouseUtility,
	// 		StorehouseStocks:  utilityStocks,
	// 	}
	// 	utilityStocksDetailed.GetTotalQuantitiesLeft()
	// 	storehouseInventoryByStocks = append(
	// 		storehouseInventoryByStocks,
	// 		utilityStocksDetailed,
	// 	)
	// }
	if filters.ShowLessQty {
		storehouseInventoryByUtility, totalPages, stocksError = interactor.StorehouseStocksDataSource.GetStorehouseInventoryUtilitiesDetailed(filters)
		data["inventory"] = storehouseInventoryByUtility
	} else {
		storehouseInventoryByStocks, totalPages, stocksError = interactor.StorehouseStocksDataSource.GetStorehouseInventoryStocks(filters)
		data["inventory"] = storehouseInventoryByStocks
	}
	if stocksError != nil {
		fmt.Println(stocksError.Error())
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        stocksError,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data:       data,
		ExtraInfo: []map[string]interface{}{
			{"pages": totalPages},
		},
	}
}

func (interactor StorehouseStocksInteractor) UpdateStorehouseStock(id string, stock models.StorehouseStock) models.Responser {
	err := interactor.StorehouseStocksDataSource.UpdateStorehouseStock(id, stock)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	stockUpdated, err := interactor.StorehouseStocksDataSource.GetStorehouseStockById(id)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"stock": stockUpdated,
		},
	}
}

func (interactor StorehouseStocksInteractor) DeleteStorehouseStock(id string) models.Responser {
	err := interactor.StorehouseStocksDataSource.DeleteStorehouseStock(id)
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
