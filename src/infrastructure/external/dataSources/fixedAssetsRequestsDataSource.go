package datasources

import (
	"cecan_inventory/domain/models"
	"errors"

	"gorm.io/gorm"
)

type FixedAssetsRequetsDataSource struct {
	DbPsql *gorm.DB
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
