package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type StorehouseUtilitiesInteractor struct {
	StorehouseUtilitiesDataSource datasources.StorehouseUtilitiesDataSource
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

func (interactor StorehouseUtilitiesInteractor) GetStorehouseUtilities() models.Responser {
	storehouseUtilities, err := interactor.StorehouseUtilitiesDataSource.GetStorehouseUtilities()
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
