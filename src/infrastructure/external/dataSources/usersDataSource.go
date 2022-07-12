package datasources

import (
	"cecan_inventory/src/infrastructure/storage/models"

	"gorm.io/gorm"
)

type UserDataSource struct {
	psql *gorm.DB
}

func (dataSrc *UserDataSource) Login() []models.User {
	var users []models.User
	dataSrc.psql.Find(&users)
	return users
}
