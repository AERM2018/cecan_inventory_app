package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

func CreateRoles(db *gorm.DB) error {
	for _, rol := range mocks.GetRolesMock() {
		err := db.FirstOrCreate(&rol, models.Role{Name: rol.Name}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
