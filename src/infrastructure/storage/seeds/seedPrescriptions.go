package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"os"

	"gorm.io/gorm"
)

func CreatePrescriptions(db *gorm.DB) error {
	if os.Getenv("GO_ENV") == "TEST" {
		var (
			doctorUser models.User
		)
		prescriptions := mocks.GetPrescriptionMockSeed()

		// Get a user in DB that has a doctor role
		db.Where("role_id = ?", mocks.GetRolesMock("medico")[0].Id).First(&doctorUser)
		for _, prescription := range prescriptions {
			var err error
			prescriptionInfo := models.Prescription{
				Id:                   prescription.Id,
				UserId:               doctorUser.Id,
				PatientName:          prescription.PatientName,
				Instructions:         prescription.Instructions,
				PrescriptionStatusId: prescription.PrescriptionStatusId,
			}
			err = db.Create(&prescriptionInfo).Error
			if err != nil {
				return err
			}
			for _, medicine := range prescription.Medicines {
				db.Create(&medicine)
			}
		}

	}
	return nil
}
