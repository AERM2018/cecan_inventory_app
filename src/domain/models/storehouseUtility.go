package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	StorehouseUtility struct {
		Key             string         `gorm:"primaryKey" json:"key"`
		GenericName     string         `json:"generic_name"`
		PresentationId  uuid.UUID      `gorm:"column:storehouse_utility_presentation_id" json:"presentation_id"`
		UnitId          uuid.UUID      `gorm:"column:storehouse_utility_unit_id" json:"unit_id"`
		CategoryId      uuid.UUID      `gorm:"column:storehouse_utility_category_id" json:"category_id"`
		QuantityPerUnit float32        `json:"quantity_per_unit"`
		Description     string         `json:"description"`
		CreatedAt       *time.Time     `json:"created_at,omitempty"`
		UpdatedAt       *time.Time     `json:"updated_at,omitempty"`
		DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty"`
	}

	StorehouseUtilityDetailed struct {
		Key             string                        `gorm:"primaryKey" json:"key"`
		GenericName     string                        `json:"generic_name"`
		UnitId          uuid.UUID                     `gorm:"column:storehouse_utility_unit_id;foreignKey:storehouse_utility_unit_id" json:"unit_id"`
		Unit            StorehouseUtilityUnit         `json:"unit"`
		PresentationId  uuid.UUID                     `gorm:"column:storehouse_utility_presentation_id;foreignKey:storehouse_utility_presentation_id" json:"presentation_id"`
		Presentation    StorehouseUtilityPresentation `json:"presentation"`
		CategoryId      uuid.UUID                     `gorm:"column:storehouse_utility_category_id;foreignKey:storehouse_utility_category_id" json:"category_id"`
		Category        StorehouseUtilityCategory     `json:"category"`
		QuantityPerUnit float32                       `json:"quantity_per_unit"`
		Description     string                        `json:"description"`
		CreatedAt       *time.Time                    `json:"created_at,omitempty"`
		UpdatedAt       *time.Time                    `json:"updated_at,omitempty"`
		DeletedAt       gorm.DeletedAt                `json:"deleted_at,omitempty"`
	}
)
