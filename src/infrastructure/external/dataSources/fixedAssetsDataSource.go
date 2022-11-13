package datasources

import (
	"cecan_inventory/domain/models"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type FixedAssetsDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc FixedAssetsDataSource) GetFixedAssets(filters models.FixedAssetFilters) ([]models.FixedAssetDetailed, error) {
	fixedAssets := make([]models.FixedAssetDetailed, 0)
	filtersJson := make(map[string]interface{})
	filterAsMap := structs.Map(filters)
	conditionString := ""
	fixedAssetFilterCounter := 0
	sqlInstance := dataSrc.DbPsql.Table("fixed_assets_detailed")
	for _, field := range structs.Fields(filters) {
		if filterAsMap[field.Name()] != "" {
			jsonTag := field.Tag("json")
			if strings.Contains(jsonTag, ",") {
				jsonTag = strings.Split(jsonTag, ",")[0]
			}
			filtersJson[jsonTag] = fmt.Sprintf("%v", filterAsMap[field.Name()])
		}
	}
	if len(maps.Keys(filtersJson)) > 0 {
		includeLogicalAndOperator := len(maps.Keys(filtersJson)) > 1
		for k, _ := range filtersJson {
			conditionString += fmt.Sprintf("%v = @%v", k, k)
			if includeLogicalAndOperator && fixedAssetFilterCounter < len(filterAsMap)-2 {
				conditionString += " AND "
			}
			fixedAssetFilterCounter -= 1
		}
		sqlInstance = sqlInstance.Where(conditionString, filtersJson)
	}

	err := sqlInstance.
		Find(&fixedAssets).
		Error
	if err != nil {
		return fixedAssets, err
	}
	return fixedAssets, nil
}

func (dataSrc FixedAssetsDataSource) GetFixedAssetByKey(key string) (models.FixedAssetDetailed, error) {
	var fixedAsset models.FixedAssetDetailed
	err := dataSrc.DbPsql.
		Table("fixed_assets_detailed").
		Where("key = ?", key).
		Take(&fixedAsset).
		Error
	if err != nil {
		return fixedAsset, err
	}
	return fixedAsset, nil

}

func (dataSrc FixedAssetsDataSource) CreateFixedAsset(fixedAsset models.FixedAsset) (string, error) {
	err := dataSrc.DbPsql.Omit("description", "brand", "model").Create(&fixedAsset).Error
	if err != nil {
		return "", err
	}
	return fixedAsset.Key, nil
}

func (dataSrc FixedAssetsDataSource) UpdateFixedAsset(key string, fixedAsset models.FixedAsset) error {
	err := dataSrc.DbPsql.
		Model(models.FixedAsset{}).
		Omit("description", "brand", "model", "director_user_id", "administrator_user_id").
		Where("key = ?", key).
		Updates(fixedAsset).
		Error
	if err != nil {
		return err
	}
	return nil
}
