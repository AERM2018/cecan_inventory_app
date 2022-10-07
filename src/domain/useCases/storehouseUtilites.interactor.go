package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type StorehouseUtilitiesInteractor struct {
	StorehouseUtilitiesDataSource datasources.StorehouseUtilitiesDataSource
	StorehouseStocksDataSource    datasources.StorehouseStocksDataSource
}

func (interactor StorehouseUtilitiesInteractor) CreateStorehouseUtility(utility models.StorehouseUtility) models.Responser {
	var (
		err               error
		storehouseUtility models.StorehouseUtilityDetailed
	)
	err = interactor.StorehouseUtilitiesDataSource.CreateStorehouseUtility(utility)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	storehouseUtility, err = interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilityByKey(utility.Key)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"storehouse_utility": storehouseUtility,
		},
	}
}

func (interactor StorehouseUtilitiesInteractor) CreateStorehouseUtilityStock(storehouseUtilityKey string, stock models.StorehouseStock) models.Responser {
	stock.StorehouseUtilityKey = storehouseUtilityKey
	stockId, err := interactor.StorehouseStocksDataSource.CreateStorehouseStock(stock)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	storehouseStock, err := interactor.StorehouseStocksDataSource.GetStorehouseStockById(stockId.String())
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"stock": storehouseStock,
		},
	}
}

func (interactor StorehouseUtilitiesInteractor) GetStorehouseUtilities(includeDeleted bool) models.Responser {

	storehouseUtilities, err := interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilities(includeDeleted)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"storehouse_utilities": storehouseUtilities,
		},
	}
}

func (interactor StorehouseUtilitiesInteractor) UpdateStorehouseUtility(key string, utility models.StorehouseUtility) models.Responser {
	utilityKeyUpdated, err := interactor.StorehouseUtilitiesDataSource.UpdateStorehouseUtility(key, utility)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	storehouseUtilityUpdated, err := interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilityByKey(utilityKeyUpdated)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"storehouse_utility": storehouseUtilityUpdated,
		},
	}
}

func (interactor StorehouseUtilitiesInteractor) ReactivateStorehouseUtility(key string) models.Responser {

	err := interactor.StorehouseUtilitiesDataSource.ReactivateStorehouseUtility(key)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	storehouseUtility, err := interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilityByKey(key)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"storehouse_utilities": storehouseUtility,
		},
	}
}

func (interactor StorehouseUtilitiesInteractor) DeleteStorehouseUtility(key string) models.Responser {
	err := interactor.StorehouseUtilitiesDataSource.DeleteStorehouseUtility(key)
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

func (interactor StorehouseUtilitiesInteractor) GetStorehouseUtilityByKey(key string) models.Responser {
	storehouseUtility, err := interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilityByKey(key)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"storehouse_utility": storehouseUtility,
		},
	}
}

func (interactor StorehouseUtilitiesInteractor) GetStorehouseUtilityCategories() models.Responser {
	storehouseUtiliyCategories, err := interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilityCategories()
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"storehouse_utility_categories": storehouseUtiliyCategories,
		},
	}
}

func (interactor StorehouseUtilitiesInteractor) GetStorehouseUtilityPresentations() models.Responser {
	storehouseUtiliyPresentations, err := interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilityPresentations()
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"storehouse_utility_presentations": storehouseUtiliyPresentations,
		},
	}
}

func (interactor StorehouseUtilitiesInteractor) GetStorehouseUtilityUnits() models.Responser {
	storehouseUtiliyUnits, err := interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilityUnits()
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"storehouse_utility_units": storehouseUtiliyUnits,
		},
	}
}
