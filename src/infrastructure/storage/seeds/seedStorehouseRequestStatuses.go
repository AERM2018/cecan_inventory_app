package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

func CreateStorehouseRequestStatuses(db *gorm.DB) error {
	for _, storehouseRequestStatus := range mocks.GetStorehouseRequestStatuses() {
		err := db.FirstOrCreate(&models.StorehouseRequestStatus{}, storehouseRequestStatus).Error
		if err != nil {
			return err
		}
	}
	return nil
}
