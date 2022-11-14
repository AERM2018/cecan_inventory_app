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
