package models

import (
	"cecan_inventory/domain/common"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Prescription struct {
		Id                   uuid.UUID      `gorm:"primaryKey;default:'uuid_generate_v4()" json:"id"`
		UserId               string         `json:"user_id"`
		PrescriptionStatusId uuid.UUID      `json:"prescription_status_id"`
		Folio                int32          `json:"folio"`
		PatientName          string         `json:"patient_name"`
		Observations         string         `json:"observations"`
		Instructions         string         `json:"instructions"`
		CreatedAt            *time.Time     `json:"created_at,omitempty"`
		UpdatedAt            *time.Time     `json:"updated_at,omitempty"`
		SuppliedAt           *time.Time     `json:"supplied_at"`
		DeletedAt            gorm.DeletedAt `json:"deleted_at,omitempty"`
	}

	PrescriptionDetialed struct {
		Id                   uuid.UUID                `gorm:"primaryKey;default:'uuid_generate_v4()" json:"id"`
		UserId               string                   `json:"user_id"`
		User                 *User                    `gorm:"foreignKey:user_id" json:"user,omitempty"`
		PrescriptionStatusId uuid.UUID                `json:"prescription_status_id"`
		PrescriptionStatus   *PrescriptionsStatues    `gorm:"foreignKey:prescription_status_id" json:"prescription_status,omitempty"`
		Medicines            []PrescriptionsMedicines `gorm:"foreignKey:prescription_id" json:"medicines,omitempty"`
		Folio                int32                    `json:"folio"`
		PatientName          string                   `json:"patient_name"`
		Observations         string                   `json:"observations"`
		Instructions         string                   `json:"instructions"`
		CreatedAt            *time.Time               `json:"created_at,omitempty"`
		UpdatedAt            *time.Time               `json:"updated_at,omitempty"`
		SuppliedAt           *time.Time               `json:"supplied_at"`
		DeletedAt            gorm.DeletedAt           `json:"deleted_at,omitempty"`
	}

	PrescriptionToComplete struct {
		Observations         string                             `json:"observations"`
		PrescriptionStatusId uuid.UUID                          `json:"prescription_status_id"`
		SuppliedAt           time.Time                          `json:"supplied_at"`
		Medicines            []PrescriptionsMedicinesToComplete `json:"medicines"`
	}
)

func (prescription *Prescription) BeforeCreate(tx *gorm.DB) (err error) {
	// Assing the next folio to the prescription
	var (
		lastPrescription   Prescription
		prescriptionStatus PrescriptionsStatues
	)
	tx.Model(&Prescription{}).Order("folio DESC").First(&lastPrescription)
	prescription.Folio = lastPrescription.Folio + 1
	// Assing default prescription status

	if prescription.PrescriptionStatusId == uuid.Nil {
		tx.Model(&PrescriptionsStatues{}).Where("name = ?", "Pendiente").First(&prescriptionStatus)
		prescription.PrescriptionStatusId = prescriptionStatus.Id
	}
	return
}

func (prescriptionDetailed *PrescriptionDetialed) FilterMedicineFromPrescription(oldMedicines []PrescriptionsMedicines) ([]PrescriptionsMedicines, []PrescriptionsMedicines) {
	// The medicines returned will be inserted to the prescriptions since they don't exists already
	// The medicines left in the prescrition medicine list are the ones that will be updated
	filteredPrescriptionMedicines := make([]PrescriptionsMedicines, 0)
	prescriptionMedicinesToInsert := make([]PrescriptionsMedicines, 0)
	prescriptionMedicinesToDelete := make([]PrescriptionsMedicines, 0)
	for _, newPrescriptionMedicine := range prescriptionDetailed.Medicines {
		isMedicine, medicine := common.FindInSlice(oldMedicines, func(i interface{}) bool {
			parsed := i.(PrescriptionsMedicines)
			return newPrescriptionMedicine.MedicineKey == parsed.MedicineKey

		})

		if !isMedicine {
			prescriptionMedicinesToInsert = append(prescriptionMedicinesToInsert, newPrescriptionMedicine)
		}
		if isMedicine && newPrescriptionMedicine.Pieces != medicine.([]PrescriptionsMedicines)[0].Pieces {
			filteredPrescriptionMedicines = append(filteredPrescriptionMedicines, newPrescriptionMedicine)
		}
	}

	for _, oldPrescriptionMedicine := range oldMedicines {
		isMedicine, _ := common.FindInSlice(prescriptionDetailed.Medicines, func(i interface{}) bool {
			parsed := i.(PrescriptionsMedicines)
			return oldPrescriptionMedicine.MedicineKey == parsed.MedicineKey

		})

		if !isMedicine {
			prescriptionMedicinesToDelete = append(prescriptionMedicinesToDelete, oldPrescriptionMedicine)
		}
	}
	prescriptionDetailed.Medicines = filteredPrescriptionMedicines
	return prescriptionMedicinesToInsert, prescriptionMedicinesToDelete
}
