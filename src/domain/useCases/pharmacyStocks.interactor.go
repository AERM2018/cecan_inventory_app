package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"fmt"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type PharmacyStockInteractor struct {
	PharmacyStocksDataSource datasources.PharmacyStocksDataSource
	MedicinesDataSource      datasources.MedicinesDataSource
}

func (interactor PharmacyStockInteractor) GetPharmacyStocks(filters models.MedicinesFilters) models.Responser {
	medicineStocksDetailed := make([]models.PharmacyStocksDetails,0)
	medicineLessQty := make([]models.PharmacyStock,0)
	data := iris.Map{}
	var (
		stocksError error
		totalPages int
	)

	// medicines, errMedicines := interactor.MedicinesDataSource.GetMedicinesCatalog(filters, filters.IncludeDeleted)
	// if errMedicines != nil {
	// 	return models.Responser{
	// 		StatusCode: iris.StatusInternalServerError,
	// 		Err:        errMedicines,
	// 	}
	// }
	// for _, medicine := range medicines {
	// 	pharmacyStocks, _ := interactor.PharmacyStocksDataSource.GetPharmacyStocksByMedicineKey(medicine.Key)
	// 	stockDetailed := models.PharmacyStocksDetailed{
	// 		Medicine: medicine,
	// 		Stocks:   pharmacyStocks,
	// 	}
	// 	stockDetailed.CountAndCategorizePieces()
	// 	medicineStocksDetailed = append(medicineStocksDetailed, stockDetailed)
	// }
	if filters.ShowLessQty {
		medicineLessQty, stocksError = interactor.PharmacyStocksDataSource.GetMedicinesWithLessStockQty(filters)
		data["inventory"] = medicineLessQty
	}else{
		// Get all stocks from pharmacy
		medicineStocksDetailed, totalPages ,stocksError = interactor.PharmacyStocksDataSource.GetPharmacyStocks(filters)
		data["inventory"] = medicineStocksDetailed
	}
	if stocksError != nil {
		fmt.Println(stocksError.Error())
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err: stocksError,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: data,
		ExtraInfo: []map[string]interface{}{
			{"pages":totalPages},
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
