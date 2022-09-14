package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type PharmacyStockInteractor struct {
	PharmacyStocksDataSource datasources.PharmacyStocksDataSource
	MedicinesDataSource      datasources.MedicinesDataSource
}

func (interactor PharmacyStockInteractor) GetPharmacyStocks() models.Responser {
	var medicineStocksDetailed []models.PharmacyStocksDetailed
	medicines, errMedicines := interactor.MedicinesDataSource.GetMedicinesCatalog(true)
	if errMedicines != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        errMedicines,
		}
	}
	for _, medicine := range medicines {
		pharmacyStocks, _ := interactor.PharmacyStocksDataSource.GetPharmacyStocksByMedicineKey(medicine.Key)
		stockDetailed := models.PharmacyStocksDetailed{
			Medicine: medicine,
			Stocks:   pharmacyStocks,
		}
		stockDetailed.CountAndCategorizePieces()
		medicineStocksDetailed = append(medicineStocksDetailed, stockDetailed)
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"inventory": medicineStocksDetailed,
		},
	}
}

func (interactor PharmacyStockInteractor) UpdatePharmacyStock(id uuid.UUID, pharmacyStock models.PharmacyStockToUpdate) models.Responser {
	err := interactor.PharmacyStocksDataSource.UpdatePharmacyStock(id, pharmacyStock)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	pharmacyStockUpdated, _ := interactor.PharmacyStocksDataSource.GetPharmacyStockById(id)
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"stock": pharmacyStockUpdated,
		},
	}
}

func (interactor PharmacyStockInteractor) DeletePharmacyStock(id uuid.UUID, isPermanent bool) models.Responser {
	err := interactor.PharmacyStocksDataSource.DeletePharmacyStock(id, isPermanent)
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
