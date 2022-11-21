package datasources

import (
	"cecan_inventory/domain/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type FixedAssetsRequetsDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc FixedAssetsRequetsDataSource) GetFixedAssetsRequests(departmentId string) ([]models.FixedAssetsRequestDetailed, error) {
	fixedAssetsRequests := make([]models.FixedAssetsRequestDetailed, 0)
	sqlInstance := dataSrc.DbPsql.
		Table("fixed_assets_requests").
		Preload("Department", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "updated_at", "deleted_at")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password", "email", "created_at", "updated_at", "deleted_at")
		})
	if departmentId != "" {
		sqlInstance = sqlInstance.Where("department_id = ?", departmentId)
	}
	err := sqlInstance.Find(&fixedAssetsRequests).
		Error
	if err != nil {
		return fixedAssetsRequests, err
	}
	return fixedAssetsRequests, nil
}

func (dataSrc FixedAssetsRequetsDataSource) GetFixedAssetsRequestById(id string) (models.FixedAssetsRequestDetailed, error) {
	var fixedAssetsRequest models.FixedAssetsRequestDetailed
	fmt.Println("id", id)
	err := dataSrc.DbPsql.
		Table("fixed_assets_requests").
		Preload("FixedAssets").
		Preload("FixedAssets.FixedAsset", func(db *gorm.DB) *gorm.DB {
			return db.Raw("select * from fixed_assets_detailed")
		}).
		Preload("Department", func(db *gorm.DB) *gorm.DB {
			return db.Table("departments").Omit("created_at", "updated_at", "deleted_at")
		}).
		Preload("Department.ResponsibleUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "surname", "full_name")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password", "email", "created_at", "updated_at", "deleted_at")
		}).
		Where("id = ?", id).
		First(&fixedAssetsRequest).
		Error
	fmt.Println("error", err)
	if err != nil {
		return fixedAssetsRequest, err
	}
	return fixedAssetsRequest, nil
}

func (dataSrc FixedAssetsRequetsDataSource) CreateFixedAssetsRequest(transactionBody func(tx *gorm.DB) error) error {
	errInTransaction := dataSrc.DbPsql.Transaction(transactionBody)
	if errInTransaction != nil {
		return errInTransaction
	}
	return nil
}

func (dataSrc FixedAssetsRequetsDataSource) DeleteFixedAssetsRequest(id string, deleteFixedAssetFunc func(fixedAssetKey string) error) error {
	errInTransaction := dataSrc.DbPsql.Transaction(func(tx *gorm.DB) error {
		fixedAssetsItemsRequest := make([]models.FixedAssetsItemsRequests, 0)
		errGetting := tx.
			Where("fixed_assets_request_id = ?", id).
			Preload("FixedAsset").
			Find(&fixedAssetsItemsRequest).
			Error
		if errGetting != nil {
			return errors.New("Ocurrió un error al obtener los elementos de material fijo de la petición.")
		}
		for _, fixedAsset := range fixedAssetsItemsRequest {
			errDeleting := deleteFixedAssetFunc(fixedAsset.FixedAsset.Key)
			if errDeleting != nil {
				return errDeleting
			}
		}
		errDeletingReq := tx.
			Where("id = ?", id).
			Delete(&models.FixedAssetsRequest{}).
			Error
		if errDeletingReq != nil {
			return errors.New("Ocurrió un error al eliminar la petición.")
		}
		return nil
	})
	if errInTransaction != nil {
		return errInTransaction
	}
	return nil
}

func (dataSrc FixedAssetsRequetsDataSource) GetSignaturesInfo() models.FixedAssetsRequestSignaturesInfo {
	director := models.User{}
	subDirector := models.User{}
	dataSrc.DbPsql.Table("users").Joins("Role").Where("\"Role\".name = ?", "Director").Take(&director)
	dataSrc.DbPsql.Table("users").Joins("Role").Where("\"Role\".name = ?", "Subdirector").Take(&subDirector)
	return models.FixedAssetsRequestSignaturesInfo{
		Director:      director,
		Administrator: subDirector,
	}
}
