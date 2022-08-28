package seeds

import (
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

func CreateAdminUser(db *gorm.DB) error {
	var (
		role models.Role
		adminUserCount int64
	)
	db.Where("name = ?","Admin").First(&role)
	db.Model(&models.User{}).Where("name = ?","CECAN ADMIN").Count(&adminUserCount)
	if(adminUserCount == 0){
		adminUser := models.User{
			RoleId: role.Id.String(),
			Password: "Qwerty*123",
			Name: "CECAN ADMIN",
			Email:"admin@cecan.com",
		}
		return db.Create(&adminUser).Error
	}
	return nil
}