package seeds

import (
	"cecan_inventory/domain/mocks"

	"gorm.io/gorm"
)

func CreateRoles(db *gorm.DB) error {
	for _, rol := range mocks.GetRolesMock("") {
		err := db.FirstOrCreate(&rol, rol).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateSuperiorRoles(db *gorm.DB) error {
	for _, rol := range mocks.GetSuperiorRolesMock("") {
		err := db.FirstOrCreate(&rol, rol).Error
		if err != nil {
			return err
		}
	}
	return nil
}
