package datasources

import (
	"cecan_inventory/domain/models"
	"errors"

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

func (dataSrc RolesDataSource) GetRoleById(id string) (models.Role, error) {
	var role models.Role
	err := dataSrc.DbPsql.Model(models.Role{}).Where("id = ?", id).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return role, errors.New("Role doesn't exist")
	}
	return role, nil
}
