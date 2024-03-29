package models

import (
	"cecan_inventory/domain/common"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	StorehouseStock struct {
		Id                       uuid.UUID           `gorm:"primaryKey;default:'uuid_generate_v4()'" json:"id"`
		StorehouseUtilityKey     string              `json:"storehouse_utility_key"`
		StorehouseUtility        StorehouseUtility   `gorm:"foreignKey:storehouse_utility_key" json:"storehouse_utility"`
		QuantityParsed           float32             `json:"quantity_parsed"`
		QuantityParsedUsed       float32             `json:"quantity_parsed_used"`
		QuantityParsedLeft       float32             `json:"quantity_parsed_left"`
		QuantityPresentation     float32             `json:"quantity_presentation"`
		QuantityPresentationUsed float32             `json:"quantity_presentation_used"`
		QuantityPresentationLeft float32             `json:"quantity_presentation_left"`
		LotNumber                string              `json:"lot_number"`
		CatalogNumber            string              `json:"catalog_number"`
		SemaforizationColor      SemaforizationColor `json:"semaforization_color"`
		ExpiresAt                time.Time           `json:"expires_at"`
		CreatedAt                *time.Time          `json:"created_at,omitempty"`
		UpdatedAt                *time.Time          `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`
	}

	StorehouseUtilityStockDetailed struct {
		Id                       uuid.UUID                 `gorm:"primaryKey;default:'uuid_generate_v4()'" json:"id"`
		StorehouseUtilityKey     string                    `json:"storehouse_utility_key"`
		StorehouseUtility        StorehouseUtilityDetailed `gorm:"foreignKey:storehouse_utility_key" json:"storehouse_utility"`
		QuantityParsed           float32                   `json:"quantity_parsed"`
		QuantityParsedUsed       float32                   `json:"quantity_parsed_used"`
		QuantityParsedLeft       float32                   `json:"quantity_parsed_left"`
		QuantityPresentation     float32                   `json:"quantity_presentation"`
		QuantityPresentationUsed float32                   `json:"quantity_presentation_used"`
		QuantityPresentationLeft float32                   `json:"quantity_presentation_left"`
		LotNumber                string                    `json:"lot_number"`
		CatalogNumber            string                    `json:"catalog_number"`
		SemaforizationColor      SemaforizationColor       `json:"semaforization_color"`
		ExpiresAt                time.Time                 `json:"expires_at"`
		CreatedAt                *time.Time                `json:"created_at,omitempty"`
		UpdatedAt                *time.Time                `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`
	}

	StorehouseUtilitsDetailedNoStocks struct {
		StorehouseUtilityKey          string                    `json:"storehouse_utility_key"`
		StorehouseUtility             StorehouseUtilityDetailed `json:"storehouse_utility"`
		TotalQuantityParsedLeft       float32                   `json:"total_quantity_parsed_left"`
		TotalQuantityPresentationLeft float32                   `json:"total_quantity_presentation_left"`
	}
	// StorehouseUtiltiyStocks struct {
	// 	Details                       StorehouseUtilityStocksDetails `json:"storehouse_utility_stocks_detailed"`
	// 	TotalQuantityParsedLeft       float32                        `json:"total:quantity_parsed_left"`
	// 	TotalQuantityPresentationLeft float32                        `json:"total_quantity_presentation_left"`
	// }
)

func (stock *StorehouseStock) BeforeCreate(tx *gorm.DB) error {
	var storehouseUtility StorehouseUtility
	err := tx.Where("key = ?", stock.StorehouseUtilityKey).First(&storehouseUtility).Error
	if err != nil {
		return err
	}
	stock.QuantityParsed = stock.QuantityPresentation * storehouseUtility.QuantityPerUnit
	stock.QuantityParsedLeft = stock.QuantityParsed
	stock.QuantityPresentationLeft = stock.QuantityPresentation
	stock.SemaforizationColor = SemaforizationColor(common.GetSemaforizationColorFromDate(stock.ExpiresAt))
	return nil
}

func (stock *StorehouseStock) BeforeUpdate(tx *gorm.DB) error {
	var storehouseUtility StorehouseUtility
	err := tx.Where("key = ?", stock.StorehouseUtilityKey).First(&storehouseUtility).Error
	if err != nil {
		return err
	}
	stock.QuantityParsed = stock.QuantityPresentation * storehouseUtility.QuantityPerUnit
	stock.QuantityParsedLeft = stock.QuantityParsed
	stock.QuantityPresentationLeft = stock.QuantityPresentation
	stock.SemaforizationColor = SemaforizationColor(common.GetSemaforizationColorFromDate(stock.ExpiresAt))
	return nil
}

// func (stocksUtilityDetailed StorehouseUtilityStockDetailed) GetTotalQuantitiesLeft() {
// 	for _, stock := range *stocksUtilityDetailed.StorehouseStocks {
// 		if stock.QuantityParsedUsed != stock.QuantityParsed {
// 			stocksUtilityDetailed.TotalQuantityParsedLeft += stock.QuantityParsed - stock.QuantityParsedUsed
// 			stocksUtilityDetailed.TotalQuantityPresentationLeft += stock.QuantityPresentation - stock.QuantityPresentationUsed
// 		}
// 	}
// }
