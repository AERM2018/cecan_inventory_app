package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	StorehouseUtility struct {
		Key                             string         `gorm:"primaryKey" json:"key"`
		GenericName                     string         `json:"generic_name"`
		StorehouseUtilityPresentationId uuid.UUID      `json:"storehouse_utility_presentation_id"`
		StorehouseUtilityUnitId         uuid.UUID      `json:"storehouse_utility_unit_id"`
		QuantityPerUnit                 float32        `json:"quantity_per_unit"`
		Description                     string         `json:"description"`
		StorehouseUtilityCategoryId     uuid.UUID      `json:"storehouse_utility_category_id"`
		CreatedAt                       *time.Time     `json:"created_at,omitempty"`
		UpdatedAt                       *time.Time     `json:"updated_at,omitempty"`
		DeletedAt                       gorm.DeletedAt `json:"deleted_at,omitempty"`
	}

	StorehouseUtilityDetailed struct {
		Key                             string                        `gorm:"primaryKey" json:"key"`
		GenericName                     string                        `json:"generic_name"`
		StorehouseUtilityUnitId         uuid.UUID                     `gorm:"foreignKey:storehouse_utility_unit_id" json:"storehouse_utility_unit_id"`
		StorehouseUtilityUnit           StorehouseUtilityUnit         `json:"storehouse_utility_unit"`
		StorehouseUtilityPresentationId uuid.UUID                     `gorm:"foreignKey:storehouse_utility_presentation_id" json:"storehouse_utility_presentation_id"`
		StorehouseUtilityPresentation   StorehouseUtilityPresentation `json:"storehouse_utility_presentation"`
		QuantityPerUnit                 float32                       `json:"quantity_per_unit"`
		Description                     string                        `json:"description"`
		StorehouseUtilityCategoryId     uuid.UUID                     `gorm:"foreignKey:storehouse_utility_category_id" json:"storehouse_utility_category_id"`
		StorehouseUtilityCategory       StorehouseUtilityCategory     `json:"storehouse_utility_category"`
		CreatedAt                       *time.Time                    `json:"created_at,omitempty"`
		UpdatedAt                       *time.Time                    `json:"updated_at,omitempty"`
		DeletedAt                       gorm.DeletedAt                `json:"deleted_at,omitempty"`
	}
)
