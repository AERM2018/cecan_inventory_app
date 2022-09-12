package datasources

import (
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

type RolesDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc RolesDataSource) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := dataSrc.DbPsql.Find(&roles).Error; err != nil {
		return roles, err
	}
	return roles, nil
}
