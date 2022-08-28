package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type PrescriptionInteractor struct {
	PrescriptionsDataSource datasources.PrescriptionsDataSource
}

func (interactor PrescriptionInteractor) CreatePrescription(prescriptionRequest models.PrescriptionDetialed) models.Responser {
	prescriptionNoMedicines := models.Prescription{
		UserId:       prescriptionRequest.UserId,
		PatientName:  prescriptionRequest.PatientName,
		Observations: prescriptionRequest.Observations,
		Instructions: prescriptionRequest.Instructions,
	}
	prescriptionId, err := interactor.PrescriptionsDataSource.CreatePrescription(prescriptionNoMedicines)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	for _, medicineForPrescription := range prescriptionRequest.Medicines {
		interactor.PrescriptionsDataSource.PutMedicineIntoPrescription(medicineForPrescription, prescriptionId)
	}
	prescriptionFound, _ := interactor.PrescriptionsDataSource.GetPrescriptionById(prescriptionId)
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"prescription": prescriptionFound,
		},
	}
}

func (interactor PrescriptionInteractor) GetPrescriptionById(id uuid.UUID) models.Responser {
	prescription, err := interactor.PrescriptionsDataSource.GetPrescriptionById(id)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data:       iris.Map{"prescription": prescription},
	}
}
