package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

func CreateStorehouseUtilitiesCategories(db *gorm.DB) error {
	for _, storehouseUtilityCategory := range mocks.GetStorehouseUtiltiesMockSeed() {
		err := db.FirstOrCreate(&models.StorehouseUtilityCategory{}, storehouseUtilityCategory).Error
		if err != nil {
			return err
		}
	}
	return nil
}
