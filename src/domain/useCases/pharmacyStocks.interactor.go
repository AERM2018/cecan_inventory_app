package usecases

import (
	"cecan_inventory/src/domain/models"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type PharmacyStockInteractor struct {
	PharmacyStocksDataSource datasources.PharmacyStocksDataSource
	MedicinesDataSource      datasources.MedicinesDataSource
}

func (interactor PharmacyStockInteractor) InsertStockOfMedicine(stock models.PharmacyStock) models.Responser {
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

func (interactor PharmacyStockInteractor) GetPharmacyStocks() models.Responser {
	var medicineStocksDetailed []models.PharmacyStocksDetailed
	medicines, errMedicines := interactor.MedicinesDataSource.GetMedicinesCatalog()
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