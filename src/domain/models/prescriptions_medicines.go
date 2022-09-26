package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrescriptionsMedicines struct {
	Id             *uuid.UUID     `gorm:"primaryKey;default:'uuid_generate_v4()'" json:"id,omitempty"`
	PrescriptionId uuid.UUID      `gorm:"foreignKey" json:"prescription_id"`
	MedicineKey    string         `gorm:"foreignKey:medicine_key" json:"medicine_key"`
	Medicine       Medicine       `json:"medicine"`
	Pieces         int16          `json:"pieces"`
	PiecesSupplied int16          `json:"pieces_supplied"`
	CreatedAt      *time.Time     `json:"created_at,omitempty"`
	UpdatedAt      *time.Time     `json:"updated_at,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
