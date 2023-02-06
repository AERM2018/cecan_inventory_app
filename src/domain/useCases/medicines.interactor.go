package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type MedicinesInteractor struct {
	MedicinesDataSource      datasources.MedicinesDataSource
	PharmacyStocksDataSource datasources.PharmacyStocksDataSource
}

func (interactor MedicinesInteractor) InsertMedicineIntoCatalog(medicine models.Medicine) models.Responser {
	creationErr := interactor.MedicinesDataSource.InsertMedicine(medicine)
	if creationErr != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        creationErr,
		}
	}
	medicineFound, _ := interactor.MedicinesDataSource.GetMedicineByKey(medicine.Key)
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"medicine": medicineFound,
		},
	}
}

func (interactor MedicinesInteractor) InsertStockOfMedicine(stock models.PharmacyStock) models.Responser {
	stockId, errOnInsertion := interactor.PharmacyStocksDataSource.InsertStockOfMedicine(stock)
	if errOnInsertion != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        errOnInsertion,
		}
	}
	newStock, errOnGetting := interactor.PharmacyStocksDataSource.GetPharmacyStockById(stockId)
	if errOnGetting != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        errOnGetting,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"stock": newStock,
		},
	}
}
func (interactor MedicinesInteractor) GetMedicinesCatalog(filters models.MedicinesFilters, includeDeleted bool) models.Responser {
	medicinesCatalog, totalPages, err := interactor.MedicinesDataSource.GetMedicinesCatalog(filters, includeDeleted)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"medicines": medicinesCatalog,
		},
		ExtraInfo: []map[string]interface{}{
			{"pages": totalPages},
		},
	}
}

func (interactor MedicinesInteractor) UpdateMedicine(key string, medicine models.Medicine) models.Responser {
	newKey, err := interactor.MedicinesDataSource.UpdateMedicine(key, medicine)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	medicineUpdated, _ := interactor.MedicinesDataSource.GetMedicineByKey(newKey)
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"medicine": medicineUpdated,
		},
	}
}

func (interactor MedicinesInteractor) DeleteMedicine(key string) models.Responser {
	err := interactor.MedicinesDataSource.DeleteMedicineByKey(key)
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

func (interactor MedicinesInteractor) ReactivateMedicine(key string) models.Responser {
	err := interactor.MedicinesDataSource.ReactivateMedicine(key)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	medicineUpdated, _ := interactor.MedicinesDataSource.GetMedicineByKey(key)
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"medicine": medicineUpdated,
		},
	}
}
