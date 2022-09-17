package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"os"

	"gorm.io/gorm"
)

func CreatePharmacyStock(db *gorm.DB) error {
	if os.Getenv("GO_ENV") == "TEST" {
		return db.FirstOrCreate(&models.PharmacyStock{}, mocks.GetPharmacyStockMockSeed()).Error
	}
	return nil
}
