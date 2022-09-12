package seeds

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

func CreateAdminUser(db *gorm.DB) error {
	var (
		role           models.Role
		adminUserCount int64
	)
	db.Where("name = ?", "Admin").First(&role)
	db.Model(&models.User{}).Where("name = ?", "CECAN ADMIN").Count(&adminUserCount)
	if adminUserCount == 0 {
		adminUser := mocks.GetUserMockSeed(role.Id.String())
		return db.Create(&adminUser).Error
	}
	return nil
}
