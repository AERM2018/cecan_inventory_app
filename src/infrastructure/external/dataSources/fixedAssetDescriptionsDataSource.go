package datasources

import (
	"cecan_inventory/domain/models"
	"database/sql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FixedAssetDescriptionDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc FixedAssetDescriptionDataSource) GetFixedAssetDescriptions(expression string) ([]models.FixedAssetDescription, error) {
	var fixedAssetDescriptions []models.FixedAssetDescription
	sqlInstance := dataSrc.DbPsql
	if expression != "" {
		sqlInstance = sqlInstance.
			Where("description ilike @query or brand ilike  @query or model ilike @query", sql.Named("query", "%"+expression+"%"))
	}
	err := sqlInstance.Find(&fixedAssetDescriptions).
		Error
	if err != nil {
		return fixedAssetDescriptions, err
	}
	return fixedAssetDescriptions, nil
}
