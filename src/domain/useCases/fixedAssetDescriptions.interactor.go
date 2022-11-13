package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
)

type FixedAssetDescriptionsInteractor struct {
	FixedAssetDescriptionsDataSource datasources.FixedAssetDescriptionDataSource
}

func (interactor FixedAssetDescriptionsInteractor) GetFixedAssetDescriptions(expression string) models.Responser {
	fixedAssetDescriptions, err := interactor.FixedAssetDescriptionsDataSource.GetFixedAssetDescriptions(expression)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"fixed_asset_descriptions": fixedAssetDescriptions,
		},
	}
}
