package seeds

import (
	"cecan_inventory/src/domain/models"

	"gorm.io/gorm"
)

func CreateRoles(db *gorm.DB, name string) error {
	return db.FirstOrCreate(&models.Role{},models.Role{Name: name}).Error
}