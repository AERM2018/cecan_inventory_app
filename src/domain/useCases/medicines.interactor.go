package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

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

func (interactor MedicinesInteractor) GetMedicinesCatalog(includeDeleted bool) models.Responser {
	medicinesCatalog, err := interactor.MedicinesDataSource.GetMedicinesCatalog(includeDeleted)
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
