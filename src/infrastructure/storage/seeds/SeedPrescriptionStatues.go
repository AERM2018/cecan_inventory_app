package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

func CreatePrescriptionStatues(db *gorm.DB) error {
	for _, prescriptionStatus := range mocks.GetPrescriptionStatuesMockSeed() {
		err := db.FirstOrCreate(&models.PrescriptionsStatues{}, models.PrescriptionsStatues{Id: prescriptionStatus.Id, Name: prescriptionStatus.Name}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
