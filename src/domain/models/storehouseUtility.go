package models

import (
	"fmt"
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
		Key               string                         `gorm:"primaryKey" json:"key"`
		GenericName       string                         `json:"generic_name"`
		UnitId            uuid.UUID                      `gorm:"column:storehouse_utility_unit_id;foreignKey:storehouse_utility_unit_id" json:"unit_id,omitempty"`
		Unit              *StorehouseUtilityUnit         `json:"unit,omitempty"`
		PresentationId    uuid.UUID                      `gorm:"column:storehouse_utility_presentation_id;foreignKey:storehouse_utility_presentation_id" json:"presentation_id"`
		Presentation      *StorehouseUtilityPresentation `json:"presentation,omitempty"`
		CategoryId        uuid.UUID                      `gorm:"column:storehouse_utility_category_id;foreignKey:storehouse_utility_category_id" json:"category_id"`
		Category          *StorehouseUtilityCategory     `json:"category,omitempty"`
		QuantityPerUnit   float32                        `json:"quantity_per_unit,omitempty"`
		FinalPresentation string                         `json:"final_presentation,omitempty"`
		Description       string                         `json:"description,omitempty"`
		CreatedAt         *time.Time                     `json:"created_at,omitempty"`
		UpdatedAt         *time.Time                     `json:"updated_at,omitempty"`
		DeletedAt         gorm.DeletedAt                 `json:"deleted_at,omitempty"`
	}

	StorehouseUtilitiesFilters struct {
		UtilityKey     string `json:"key,omitempty" json2:"utility_key"`
		UtilityName    string `json:"generic_name,omitempty" json2:"utility_name"`
		Limit          int    `json:"limit,omitempty"`
		Page           int    `json:"page,omitempty"`
		IncludeDeleted bool   `json:"include_deleted,omitempty"`
		ShowLessQty    bool   `json:"show_less_qty,omitempty"`
	}
)

func (stock *StorehouseUtilityDetailed) AfterFind(tx *gorm.DB) error {
	presentationName := ""
	unitsName := ""
	tx.Select("name").Table("storehouse_utility_presentations").Where("id = ?", stock.PresentationId).Take(&presentationName)
	tx.Select("name").Table("storehouse_utility_units").Where("id = ?", stock.UnitId).Take(&unitsName)
	stock.FinalPresentation = fmt.Sprintf("%v de %v %v", presentationName, stock.QuantityPerUnit, unitsName)
	return nil
}

func (StorehouseUtilityDetailed) TableName() string {
	return "storehouse_utilities"
}
