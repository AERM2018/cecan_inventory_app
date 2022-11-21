package usecases

import (
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type FixedAssetsRequestsInteractor struct {
	FixedAssetsRequestsDataSource datasources.FixedAssetsRequetsDataSource
	FixedAssetsDataSource         datasources.FixedAssetsDataSource
}

func (interactor FixedAssetsRequestsInteractor) GetFixedAssetsRequests(departmentId string) models.Responser {
	fixedAssetsRequests, err := interactor.FixedAssetsRequestsDataSource.GetFixedAssetsRequests(departmentId)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Data: iris.Map{
			"fixed_assets_requests": fixedAssetsRequests,
		},
	}
}

func (interactor FixedAssetsRequestsInteractor) GetFixedAssetsRequestById(id string, isPdf bool) models.Responser {
	fixedAssetsRequest, err := interactor.FixedAssetsRequestsDataSource.GetFixedAssetsRequestById(id)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusInternalServerError,
			Err:        err,
		}
	}
	if isPdf {
		fixedAssetsRequestDoc := models.FixedAssetsRequestDoc{FixedAssetRequest: fixedAssetsRequest}
		filePath, errInPdf := fixedAssetsRequestDoc.CreateDoc()
		if errInPdf != nil {
			return models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    "Ocurrió un error al generar el documento digital, intentelo mas tarde.",
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
			"fixed_assets_request": fixedAssetsRequest,
		},
	}
}

func (interactor FixedAssetsRequestsInteractor) CreateFixedAssetsRequest(
	fixedAssetsRequest models.FixedAssetsRequestDetailed,
	createAssetFunc func(fixedAsset models.FixedAsset) models.Responser,
	deleteAssetFunc func(fixedAssetKey string) models.Responser,
) models.Responser {
	requestNoAssets := models.FixedAssetsRequest{
		UserId:       fixedAssetsRequest.UserId,
		DepartmentId: fixedAssetsRequest.DepartmentId,
	}
	err := interactor.FixedAssetsRequestsDataSource.CreateFixedAssetsRequest(
		// Init a transaction making use of the data source of fixed assets
		func(tx *gorm.DB) error {
			// Create fixed asset request instance
			errCreatingReq := tx.Create(&requestNoAssets).Error
			if errCreatingReq != nil {
				return errCreatingReq
			}
			// Create fixed asset instances
			// Use create asset function from the interactor to include the logic of description reusability
			for _, asset := range fixedAssetsRequest.FixedAssets {
				res := createAssetFunc(asset.FixedAsset)
				err := res.Err
				if err != nil {
					return err
				}
				fixedAssetItemRequest := models.FixedAssetsItemsRequests{
					FixedAssetId:         maps.Values(res.Data)[0].(models.FixedAssetDetailed).Id,
					FixedAssetsRequestId: requestNoAssets.Id,
				}
				// Associate the fixed asset with the request
				errAssociating := tx.Create(&fixedAssetItemRequest).Error
				if errAssociating != nil {
					deleteAssetFunc(asset.FixedAsset.Key)
					return errAssociating
				}
			}
			return nil
		})
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Err:        err,
			Message:    err.Error(),
		}
	}
	fixedAssetsRequestCreated, _ := interactor.FixedAssetsRequestsDataSource.GetFixedAssetsRequestById(requestNoAssets.Id.String())
	return models.Responser{
		StatusCode: iris.StatusCreated,
		Data: iris.Map{
			"fixed_assets_request": fixedAssetsRequestCreated,
		},
	}
}

func (interactor FixedAssetsRequestsInteractor) DeleteFixedAssetsRequest(id string) models.Responser {
	err := interactor.FixedAssetsRequestsDataSource.DeleteFixedAssetsRequest(id, interactor.FixedAssetsDataSource.DeleteFixedAsset)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Err:        err,
			Message:    err.Error(),
		}
	}
	return models.Responser{
		StatusCode: iris.StatusNoContent,
	}
}