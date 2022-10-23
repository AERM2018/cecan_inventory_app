package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

func CreateStorehouseUtilityUnits(db *gorm.DB) error {
	for _, storehouseUtilityUnit := range mocks.GetStorehouseUtiltyUnitsMockSeed() {
		err := db.FirstOrCreate(&models.StorehouseUtilityUnit{}, storehouseUtilityUnit).Error
		if err != nil {
			return err
		}
	}
	return nil
}
