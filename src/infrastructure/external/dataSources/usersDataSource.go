package datasources

import (
	"cecan_inventory/src/infrastructure/storage/models"

	"gorm.io/gorm"
)

type UserDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc *UserDataSource) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	res := dataSrc.DbPsql.Where(&models.User{Email: email}).First(&user)
	if res.RowsAffected < 1 {
		return user, res.Error;
	}
	return user, nil
}
