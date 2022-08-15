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
		PrescriptionStatus   *Prescriptions_statues   `gorm:"foreignKey:prescription_status_id" json:"prescription_status,omitempty"`
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

	PrescriptionsMedicinesReq struct {
		MedicineKey string `json:"medicine_key"`
		Pieces      int16  `json:"pieces"`
	}
)

func (prescription *Prescription) BeforeCreate(tx *gorm.DB) (err error) {
	// Assing the next folio to the prescription
	var (
		lastPrescription   Prescription
		prescriptionStatus Prescriptions_statues
	)
	tx.Model(&Prescription{}).Order("folio DESC").First(&lastPrescription)
	prescription.Folio = lastPrescription.Folio + 1
	// Assing default prescription status

	tx.Model(&Prescriptions_statues{}).Where("name = ?", "Pendiente").First(&prescriptionStatus)
	prescription.PrescriptionStatusId = prescriptionStatus.Id
	return
}

// func (prescriptionDetailed PrescriptionDetialed) ToJSON() (string, error) {
// 	medicines := make()
// }
