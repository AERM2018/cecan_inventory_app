package usecases

import (
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type FixedAssetsInteractor struct {
	FixedAssetsDataSource            datasources.FixedAssetsDataSource
	FixedAssetDescriptionsDataSource datasources.FixedAssetDescriptionDataSource
	DepartmentsDataSource            datasources.DepartmentDataSource
	UserDataSource                   datasources.UserDataSource
}

func (interactor FixedAssetsInteractor) GetFixedAssets(filters models.FixedAssetFilters, datesDelimiter []string, isPdf bool, page int, limit int, offset int) models.Responser {
	fixedAssets, numPages, err := interactor.FixedAssetsDataSource.GetFixedAssets(filters, datesDelimiter, page, limit, offset)
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
		ExtraInfo: []map[string]interface{}{
			{"pages": numPages},
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
	var (
		descriptionId uuid.UUID
		departmentId  uuid.UUID
		department    models.DepartmentDetailed
		errNotFound   error
	)
	data, _ := common.ReadDataFromCsv(filePath)
	isDescription := false
	superiorUsers, _ := interactor.UserDataSource.GetSuperiorUsers()
	for _, row := range data {
		// Donde se encuentra la información en la fila
		// Clave esta en el indice 0
		// Descripcion esta en el indice 1
		// Marca esta en el indice 2
		// Modelo esta en el indice 3
		_, errNotFound = interactor.FixedAssetsDataSource.GetFixedAssetByKey(row[0])
		if errNotFound != nil {

			descriptionModel := models.FixedAssetDescription{
				Description: strings.ToUpper(row[1]),
				Brand:       strings.ToUpper(row[2]),
				Model:       strings.ToUpper(row[3]),
			}
			_, descriptionId = interactor.FixedAssetDescriptionsDataSource.GetSimilarFixedAssetDescriptions(descriptionModel)
			// fmt.Printf("%v is in db %v\n", descriptionId, isDescription)
			if !isDescription {
				descriptionId, _ = interactor.FixedAssetDescriptionsDataSource.CreateFixedAssetDescription(descriptionModel)
			}
			department, errNotFound = interactor.DepartmentsDataSource.GetDepartmentByName(row[8])
			// fmt.Printf("%v is in db\n", department.Id)
			if errNotFound != nil {
				// create department if it doesn't exist
				departmentId, _ = interactor.DepartmentsDataSource.CreateDepartment(models.Department{
					FloorNumber:       row[10],
					Name:              row[8],
					ResponsibleUserId: nil,
				})
				department, _ = interactor.DepartmentsDataSource.GetDepartmentById(departmentId.String())
			} else {
				departmentId = department.Id
			}
			// fmt.Printf("descriptionId %v", descriptionId)
			fixedAsset := models.FixedAsset{
				Key:                         row[0],
				Description:                 row[1],
				Brand:                       row[2],
				Model:                       row[3],
				FixedAssetDescriptionId:     &descriptionId,
				DepartmentId:                departmentId,
				DepartmentResponsibleUserId: department.ResponsibleUser.Id,
				Series:                      row[4],
				Type:                        row[5],
				PhysicState:                 row[6],
				Observation:                 row[9],
				DirectorUserId:              superiorUsers[0].Id,
				AdministratorUserId:         superiorUsers[1].Id,
			}
			fmt.Printf("%v is in db\n", fixedAsset.FixedAssetDescriptionId)
			fixedAssetId, err := interactor.FixedAssetsDataSource.CreateFixedAsset(fixedAsset)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Print("Fixed asset id: ")
			fmt.Println(fixedAssetId)

		}

	}
	// interactor.FixedAssetsDataSource.CreateFixedAsset(models.FixedAssetDescription{
	// 	Description: ,
	// })
	return models.Responser{
		StatusCode: iris.StatusNoContent,
	}
}
