package datasources

import (
	"cecan_inventory/domain/models"
	"errors"

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
	err := dataSrc.DbPsql.
		Table("fixed_assets_requests").
		Preload("FixedAssets").
		Preload("FixedAssets.FixedAsset", func(db *gorm.DB) *gorm.DB {
			return db.Raw("select * from fixed_assets_detailed")
		}).
		Preload("Department", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "updated_at", "deleted_at")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password", "email", "created_at", "updated_at", "deleted_at")
		}).
		Where("id = ?", id).
		Find(&fixedAssetsRequest).
		Error
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
			errDeleting := deleteFixedAssetFunc(fixedAsset.FixedAssetKey)
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
