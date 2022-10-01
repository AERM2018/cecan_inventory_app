package seeds

import (
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

func CreatePrescriptionStatues(db *gorm.DB) error {
	prescriptionStatusNames := []string{"Pendiente", "Completada"}
	for _, statusName := range prescriptionStatusNames {
		err := db.Create(&models.PrescriptionsStatues{Name: statusName}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
