package usecases

import (
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"strings"

	"github.com/kataras/iris/v12"
)

type FixedAssetsInteractor struct {
	FixedAssetsDataSource            datasources.FixedAssetsDataSource
	FixedAssetDescriptionsDataSource datasources.FixedAssetDescriptionDataSource
}

func (interactor FixedAssetsInteractor) GetFixedAssets(filters models.FixedAssetFilters, datesDelimiter []string, isPdf bool) models.Responser {
	fixedAssets, err := interactor.FixedAssetsDataSource.GetFixedAssets(filters, datesDelimiter)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	if isPdf {
		fixedAssetsReport := models.FixedAssetsReportDoc{FixedAssets: fixedAssets}
		filePath, errInPdf := fixedAssetsReport.CreateDoc()
		if errInPdf != nil {
			return models.Responser{
				StatusCode: iris.StatusInternalServerError,
				Err:        errInPdf,
			}
		}

		return models.Responser{
			ExtraInfo: []map[string]interface{}{
				{"file": filePath},
			},
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"fixed_assets": fixedAssets,
		},
	}
}

func (interactor FixedAssetsInteractor) CreateFixedAsset(fixedAsset models.FixedAsset) models.Responser {
	fixedAssetDescription := models.FixedAssetDescription{
		Description: strings.ToUpper(fixedAsset.Description),
		Brand:       strings.ToUpper(fixedAsset.Brand),
		Model:       strings.ToUpper(fixedAsset.Model),
	}
	isSimilarFound, similarDescriptionId := interactor.FixedAssetDescriptionsDataSource.GetSimilarFixedAssetDescriptions(fixedAssetDescription)
	if !isSimilarFound {
		fixedAssetDescription.Id, _ = interactor.FixedAssetDescriptionsDataSource.CreateFixedAssetDescription(fixedAssetDescription)
	} else {
		fixedAssetDescription.Id = similarDescriptionId
	}
	fixedAsset.FixedAssetDescriptionId = &fixedAssetDescription.Id
	_, err := interactor.FixedAssetsDataSource.CreateFixedAsset(fixedAsset)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	fixedAssetFound, _ := interactor.FixedAssetsDataSource.GetFixedAssetByKey(fixedAsset.Key)
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"fixed_asset": fixedAssetFound,
		},
	}
}

func (interactor FixedAssetsInteractor) UpdateFixedAsset(key string, fixedAsset models.FixedAsset) models.Responser {
	fixedAssetDescription := models.FixedAssetDescription{
		Description: strings.ToUpper(fixedAsset.Description),
		Brand:       strings.ToUpper(fixedAsset.Brand),
		Model:       strings.ToUpper(fixedAsset.Model),
	}
	// Get data previously stored before updating
	fixAssetFound, err := interactor.FixedAssetsDataSource.GetFixedAssetByKey(key)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	// Verify if the description changed so as to look for a description similarity
	if !(fixedAsset.Description == fixAssetFound.Description &&
		fixedAsset.Brand == fixAssetFound.Brand &&
		fixedAsset.Model == fixAssetFound.Model) {
		isSimilarFound, similarDescriptionId := interactor.FixedAssetDescriptionsDataSource.GetSimilarFixedAssetDescriptions(fixedAssetDescription)
		if !isSimilarFound {
			fixedAssetDescription.Id, _ = interactor.FixedAssetDescriptionsDataSource.CreateFixedAssetDescription(fixedAssetDescription)
		} else {
			fixedAssetDescription.Id = similarDescriptionId
		}
		fixedAsset.FixedAssetDescriptionId = &fixedAssetDescription.Id
	}
	interactor.FixedAssetsDataSource.UpdateFixedAsset(key, fixedAsset)
	// Look for the instance updated
	// Use the key of the req data in case it would have changed
	fixedAssetUpdated, errGettingUpdated := interactor.FixedAssetsDataSource.GetFixedAssetByKey(fixedAsset.Key)
	if errGettingUpdated != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        errGettingUpdated,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"fixed_asset": fixedAssetUpdated,
		},
	}
}

func (interactor FixedAssetsInteractor) DeleteFixedAsset(key string) models.Responser {
	err := interactor.FixedAssetsDataSource.DeleteFixedAsset(key)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusNoContent,
	}
}

func (interactor FixedAssetsInteractor) UploadFileDataToDb(filePath string) models.Responser {
	common.UploadCsvToDb(filePath)
	return models.Responser{
		StatusCode: iris.StatusNoContent,
	}
}
