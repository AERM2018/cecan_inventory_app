package usecases

import (
	"cecan_inventory/src/domain/models"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type MedicinesInteractor struct {
	MedicinesDataSource datasources.MedicinesDataSource
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

func (interactor MedicinesInteractor) GetMedicinesCatalog() models.Responser {
	medicinesCatalog, err := interactor.MedicinesDataSource.GetMedicinesCatalog()
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"medicines": medicinesCatalog,
		},
	}
}
