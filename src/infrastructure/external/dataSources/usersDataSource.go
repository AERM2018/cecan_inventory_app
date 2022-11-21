package datasources

import (
	"cecan_inventory/domain/models"
	"errors"

	"gorm.io/gorm"
)

type UserDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc UserDataSource) GetUserByEmailOrId(username string) (models.User, error) {
	var user models.User
	res := dataSrc.DbPsql.
		Joins("Role").
		Omit("updated_at", "deleted_at").
		Where("email = ? or \"users\".id = ?", username, username).First(&user)
	if res.RowsAffected < 1 {
		return user, res.Error
	}
	return user, nil
}

func (dataSrc UserDataSource) CreateUser(user models.User) (models.User, error) {
	var newUserOrFound models.User
	res := dataSrc.DbPsql.FirstOrCreate(&newUserOrFound, &user)
	if res.RowsAffected < 1 {
		if res.Error != nil {
			return newUserOrFound, res.Error
		}
		return newUserOrFound, errors.New("Un usuario con el email " + user.Email + " ya existe.")
	}
	return newUserOrFound, nil
}
