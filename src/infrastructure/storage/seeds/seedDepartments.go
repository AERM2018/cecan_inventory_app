package seeds

import (
	"cecan_inventory/domain/mocks"

	"gorm.io/gorm"
)

func CreateDepartments(db *gorm.DB) error {
	departments := mocks.GetDepartmentsMockSeed()
	for _, department := range departments {
		db.FirstOrCreate(&department, department)
	}
	return nil
}
