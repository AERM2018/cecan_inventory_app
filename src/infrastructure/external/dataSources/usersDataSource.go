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

func (dataSrc UserDataSource) GetUsers() ([]models.User, error) {
	users := make([]models.User, 0)
	err := dataSrc.DbPsql.
		Omit("password").
		Find(&users).
		Error
	if err != nil {
		return users, err
	}
	return users, nil
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

func (dataSrc UserDataSource) UpdateUserPassword(user models.User) error {
	err := dataSrc.DbPsql.
		Model(&models.User{}).
		Where("id = ?", user.Id).
		Update("password", user.Password).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc UserDataSource) GetSuperiorUsers() ([]models.User, error) {
	users := make([]models.User, 0)
	err := dataSrc.DbPsql.
		Omit("password").
		Joins("Role").
		Where("\"Role\".name = ? or \"Role\".name = ?", "Director", "Subdirector").
		Find(&users).
		Error
	if err != nil {
		return users, err
	}
	return users, nil
}
