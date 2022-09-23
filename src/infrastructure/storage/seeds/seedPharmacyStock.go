package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"os"

	"gorm.io/gorm"
)

func CreatePharmacyStock(db *gorm.DB) error {
	if os.Getenv("GO_ENV") == "TEST" {
		pharmacyStocks := mocks.GetPharmacyStockMockSeed()
		for _, stock := range pharmacyStocks {
			if err := db.FirstOrCreate(&models.PharmacyStock{}, stock).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
