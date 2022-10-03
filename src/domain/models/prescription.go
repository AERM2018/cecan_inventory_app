package models

import (
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
