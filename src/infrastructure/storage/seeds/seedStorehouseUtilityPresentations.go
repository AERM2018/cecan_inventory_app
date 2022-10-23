package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

func CreateStorehouseUtilityPresentations(db *gorm.DB) error {
	for _, storehouseUtilityPresentation := range mocks.GetStorehouseUtiltyPresentationsMockSeed() {
		err := db.FirstOrCreate(&models.StorehouseUtilityPresentation{}, storehouseUtilityPresentation).Error
		if err != nil {
			return err
		}
	}
	return nil
}
