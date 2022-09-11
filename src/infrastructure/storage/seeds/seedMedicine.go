package seeds

import (
	"cecan_inventory/domain/mocks"
	"os"

	"gorm.io/gorm"
)

func CreateMedicines(db *gorm.DB) error {
	if os.Getenv("GO_ENV") == "TEST" {
		medicines := mocks.GetMedicineMockSeed()
		for _, medicine := range medicines {
			db.Create(&medicine)
		}
	}

	return nil
}
