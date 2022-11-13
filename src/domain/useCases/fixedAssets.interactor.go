package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"strings"

	"github.com/kataras/iris/v12"
)

type FixedAssetsInteractor struct {
	FixedAssetsDataSource            datasources.FixedAssetsDataSource
	FixedAssetDescriptionsDataSource datasources.FixedAssetDescriptionDataSource
}

func (interactor FixedAssetsInteractor) GetFixedAssets(filters models.FixedAssetFilters) models.Responser {
	fixedAssets, err := interactor.FixedAssetsDataSource.GetFixedAssets(filters)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"fixed_assets": fixedAssets,
		},
	}
}
