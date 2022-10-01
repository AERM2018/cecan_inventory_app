package seeds

import (
	"cecan_inventory/domain/mocks"

	"gorm.io/gorm"
)

func CreatePrescriptions(db *gorm.DB) error {
	prescriptions := mocks.GetPrescriptionMockSeed()
	for _, prescription := range prescriptions {
		err := db.Create(&prescription).Error
		if err != nil {
			return err
		}
	}
	return nil
}
