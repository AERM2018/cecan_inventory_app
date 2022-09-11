package seeds

import (
	"cecan_inventory/domain/mocks"

	"gorm.io/gorm"
)

func CreateRoles(db *gorm.DB) error {
	for _, rol := range mocks.GetRolesMock() {
		err := db.Create(&rol).Error
		if err != nil {
			return err
		}
	}
	return nil
}
