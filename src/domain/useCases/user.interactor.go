package usecases

import (
	"cecan_inventory/src/domain/models"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"
	"errors"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type UserInteractor struct {
	UserDataSource datasources.UserDataSource
}

func (interacor UserInteractor) LoginUser(credentials models.AccessCredentials) models.Responser{
		user, err := interacor.UserDataSource.GetUserByEmail(credentials.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Responser{
				StatusCode:iris.StatusNotFound,
				Message:"Invalid credentials.",
				Err:nil,
				Data:nil,
			}
		}
	}
	if isCorrectPassword := user.CheckPassword(credentials.Password); !isCorrectPassword {
		return models.Responser{
				StatusCode:iris.StatusNotFound,
				Message:"Invalid credentials.",
				Err:nil,
				Data:nil,
			}
		}
		return models.Responser{
					StatusCode:iris.StatusOK,
					Message:"",
					Err:nil,
					Data:iris.Map{"user": user.ToJSON()},
				}
	}

func (iterator UserInteractor) SignUpUser(user models.User) models.Responser{
	newUserRecord, err := iterator.UserDataSource.CreateUser(user)
	if err != nil{
		return models.Responser{
				StatusCode:iris.StatusBadRequest,
				Message:err.Error(),
				Err:nil,
				Data:nil,
			}
	}
	return models.Responser{
					StatusCode:iris.StatusOK,
					Message:"",
					Err:nil,
					Data:iris.Map{"user": newUserRecord.ToJSON()},
				}
}