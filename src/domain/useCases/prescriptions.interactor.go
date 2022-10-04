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
	prescriptionId, err := interactor.PrescriptionsDataSource.CreatePrescription(prescriptionNoMedicines, prescriptionRequest.Medicines)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Err:        err,
			Message:    err.Error(),
		}
	}
	prescriptionFound, _ := interactor.PrescriptionsDataSource.GetPrescriptionById(prescriptionId)
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"prescription": prescriptionFound,
		},
	}
}

func (interactor PrescriptionInteractor) GetPrescriptions(userId string) models.Responser {
	prescriptions, err := interactor.PrescriptionsDataSource.GetPrescriptions(userId)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data:       iris.Map{"prescriptions": prescriptions},
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

func (interactor PrescriptionInteractor) UpdatePrescription(id string, prescriptionRequest models.PrescriptionDetialed) models.Responser {
	var oldPrescription models.PrescriptionDetialed
	prescriptionNoMedicines := models.Prescription{
		UserId:       prescriptionRequest.UserId,
		PatientName:  prescriptionRequest.PatientName,
		Observations: prescriptionRequest.Observations,
		Instructions: prescriptionRequest.Instructions,
	}
	uuidFromString, _ := uuid.Parse(id)
	// Get prescription to compare medicines list to the one received
	oldPrescription, _ = interactor.PrescriptionsDataSource.GetPrescriptionById(uuidFromString)
	// Get prescription medicine list differences (INSERT AND DELETE)
	medicinesToInsert, medicinesToDelete := prescriptionRequest.FilterMedicineFromPrescription(oldPrescription.Medicines)
	for _, prescriptionMedicine := range medicinesToInsert {
		interactor.PrescriptionsDataSource.PutMedicineIntoPrescription(prescriptionMedicine, uuidFromString)
	}
	for _, prescriptionMedicine := range medicinesToDelete {
		interactor.PrescriptionsDataSource.RemoveMedicineFromPrescription(prescriptionMedicine.MedicineKey, oldPrescription.Id)
	}
	// Update prescription info and medicines quantities which had changed
	err := interactor.PrescriptionsDataSource.UpdatePrescription(id, prescriptionNoMedicines, prescriptionRequest.Medicines)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Err:        err,
			Message:    err.Error(),
		}
	}
	// Consult prescription updated
	prescriptionFound, _ := interactor.PrescriptionsDataSource.GetPrescriptionById(uuidFromString)
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"prescription": prescriptionFound,
		},
	}
}

func (interactor PrescriptionInteractor) CompletePrescription(id string, prescription models.PrescriptionToComplete) models.Responser {
	err := interactor.PrescriptionsDataSource.CompletePrescription(id, prescription)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Err:        err,
			Message:    err.Error(),
		}
	}
	uuidFromString, _ := uuid.Parse(id)
	prescriptionFound, _ := interactor.PrescriptionsDataSource.GetPrescriptionById(uuidFromString)
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"prescription": prescriptionFound,
		},
	}
}

func (interactor PrescriptionInteractor) DeletePrescription(id string) models.Responser {
	err := interactor.PrescriptionsDataSource.DeletePrescription(id)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Err:        err,
			Message:    err.Error(),
		}
	}
	return models.Responser{
		StatusCode: iris.StatusNoContent,
	}
}
