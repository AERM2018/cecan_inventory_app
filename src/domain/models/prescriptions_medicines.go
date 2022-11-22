package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PrescriptionsMedicines struct {
		Id                 *uuid.UUID     `gorm:"primaryKey;default:'uuid_generate_v4()'" json:"id,omitempty"`
		PrescriptionId     uuid.UUID      `gorm:"foreignKey" json:"prescription_id"`
		MedicineKey        string         `gorm:"foreignKey:medicine_key" json:"medicine_key"`
		Medicine           Medicine       `json:"details"`
		Pieces             int16          `json:"pieces"`
		PiecesSupplied     int16          `json:"pieces_supplied"`
		LastPiecesSupplied int16          `json:"last_pieces_supplied"`
		CreatedAt          *time.Time     `json:"created_at,omitempty"`
		UpdatedAt          *time.Time     `json:"updated_at,omitempty"`
		DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	}

	PrescriptionsMedicinesToComplete struct {
		MedicineKey    string `json:"medicine_key"`
		PiecesSupplied int16  `json:"pieces_supplied"`
	}
)

func (prescriptionMedicines PrescriptionsMedicines) BeforeCreate(tx *gorm.DB) error {
	errNotFound := tx.Where("key = ?", prescriptionMedicines.MedicineKey).First(&Medicine{}).Error
	if errors.Is(errNotFound, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("No fue posible crear la receta, el medicamento con clave: %v no existe", prescriptionMedicines.MedicineKey))
	}
	return nil
}
