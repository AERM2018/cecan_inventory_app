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

func (dataSrc FixedAssetDescriptionDataSource) GetSimilarFixedAssetDescriptions(fixedAssetDescription models.FixedAssetDescription) (bool, uuid.UUID) {
	similarFixedAssetDescriptions := make([]models.FixedAssetDescription, 0)
	dataSrc.DbPsql.
		Where("upper(description) = ? AND upper(brand) = ? AND upper(model) = ?", fixedAssetDescription.Description, fixedAssetDescription.Brand, fixedAssetDescription.Model).
		Find(&similarFixedAssetDescriptions)
	if len(similarFixedAssetDescriptions) == 0 {
		return false, uuid.Nil
	}
	return true, similarFixedAssetDescriptions[0].Id

}

func (dataSrc FixedAssetDescriptionDataSource) CreateFixedAssetDescription(fixedAssetDescription models.FixedAssetDescription) (uuid.UUID, error) {
	err := dataSrc.DbPsql.Create(&fixedAssetDescription).Error
	if err != nil {
		return uuid.Nil, err
	}
	return fixedAssetDescription.Id, nil
}
