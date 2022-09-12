package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type RolesInteractor struct {
	RolesDataSource datasources.RolesDataSource
}

func (interactor RolesInteractor) GetRoles() models.Responser {
	roles, err := interactor.RolesDataSource.GetRoles()
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"roles": roles,
		},
	}
}
