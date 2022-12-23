package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type UsersInteractor struct {
	UserDataSource datasources.UserDataSource
}

func (interactor UsersInteractor) GetUsers() models.Responser {
	users, err := interactor.UserDataSource.GetUsers()
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"users": users,
		},
	}
}
