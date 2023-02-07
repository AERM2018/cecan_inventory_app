package datasources

import (
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	"errors"
	"fmt"
	"math"

	"github.com/fatih/structs"
	"gorm.io/gorm"
)

type FixedAssetsDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc FixedAssetsDataSource) GetFixedAssets(filters models.FixedAssetFilters, datesDelimiter []string, page int, limit int, offset int) ([]models.FixedAssetDetailed, float64, error) {
	var totalRecordsCounter int64
	fixedAssets := make([]models.FixedAssetDetailed, 0)
	// filtersJson := make(map[string]interface{})
	// filterAsMap := structs.Map(filters)
	// conditionString := ""
	// fixedAssetFilterCounter := 0
	sqlInstance := dataSrc.DbPsql.Table("fixed_assets_detailed")
	// Convert struct property names to json tags and remove the ones which are empty
	// for _, field := range structs.Fields(filters) {
	// 	if filterAsMap[field.Name()] != "" {
	// 		jsonTag := field.Tag("json")
	// 		if strings.Contains(jsonTag, ",") {
	// 			jsonTag = strings.Split(jsonTag, ",")[0]
	// 		}
	// 		filtersJson[jsonTag] = fmt.Sprintf("%v", filterAsMap[field.Name()])
	// 	}
	// }
	// if len(maps.Keys(filtersJson)) > 0 {
	// 	keys := reflect.ValueOf(filtersJson).MapKeys()
	// 	fmt.Println(keys)
	// 	fmt.Println(filtersJson[keys[0].String()])

	// 	includeLogicalAndOperator := len(maps.Keys(filtersJson)) > 1
	// 	for _,k := range keys {
	// 		conditionString += fmt.Sprintf("%v LIKE %v%v%v", k, "'%",filtersJson[k.String()],"%'")
	// 		if includeLogicalAndOperator && fixedAssetFilterCounter+1 < len(filtersJson) {
	// 			conditionString += " OR "
	// 		}
	// 		fixedAssetFilterCounter += 1
	// 	}
	// 	conditionString += fmt.Sprintf(" AND \"created_at\" BETWEEN %v AND %v", datesDelimiter[0], datesDelimiter[1])
	// 	sqlInstance = sqlInstance.Where(conditionString)
	// } else {
	// 	conditionString += fmt.Sprintf("\"created_at\" BETWEEN %v AND %v", datesDelimiter[0], datesDelimiter[1])
	// 	sqlInstance = sqlInstance.Where(conditionString)
	// }
	conditionString := common.StructJsonSerializer(structs.Map(filters), "json")
	if conditionString != "" {
		conditionString += " AND "
	}
	conditionString += fmt.Sprintf("\"created_at\" BETWEEN %v AND %v", datesDelimiter[0], datesDelimiter[1])
	sqlInstance.Count(&totalRecordsCounter)
	err := sqlInstance.
		Offset((page - 1) * offset).
		Limit(limit).
		Find(&fixedAssets).
		Error
	if err != nil {
		return fixedAssets, 0, err
	}
	totalPages := math.Round(float64(int(totalRecordsCounter) / limit))
	if totalPages*float64(limit) < float64(totalRecordsCounter) {
		totalPages += 1
	}
	return fixedAssets, totalPages, nil
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
		Table("fixed_assets").
		Select("key", "series", "type", "physic_state", "observation").
		Where("key = ?", key).
		Updates(fixedAsset).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc FixedAssetsDataSource) DeleteFixedAsset(key string) error {
	errInTransaction := dataSrc.DbPsql.Transaction(func(tx *gorm.DB) error {
		errDisassociating := tx.
			Where("fixed_asset_id in(select id from fixed_assets where key = ?)", key).
			Delete(&models.FixedAssetsItemsRequests{}).
			Error
		if errDisassociating != nil {
			return errors.New("Ocurrió un error al remover el material de activo fijo de la petición")
		}
		errDeleting := tx.
			Where("key = ?", key).
			Delete(&models.FixedAsset{}).
			Error
		if errDeleting != nil {
			return errors.New("Ocurrió un error al eliminar el material de activo fijo")
		}

		return nil
	})
	if errInTransaction != nil {
		return errInTransaction
	}
	return nil
}

func (dataSrc FixedAssetsDataSource) IsFixedAssetWithSeries(series string, key string) bool {
	var fixedAsset models.FixedAsset
	whereValues := []interface{}{fmt.Sprintf("'%v'", series)}
	whereCondition := "\"series\" = ?"
	if key != "" {
		whereCondition += " AND \"key\" != ?"
		whereValues = append(whereValues, fmt.Sprintf("'%v'", key))
	}
	err := dataSrc.DbPsql.
		Where(whereCondition, whereValues...).
		Take(&fixedAsset).
		Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
