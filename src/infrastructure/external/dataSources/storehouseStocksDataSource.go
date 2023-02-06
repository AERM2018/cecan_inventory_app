package datasources

import (
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorehouseStocksDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc StorehouseStocksDataSource) CreateStorehouseStock(stock models.StorehouseStock) (uuid.UUID, error) {
	err := dataSrc.DbPsql.Select(
		"storehouse_utility_key",
		"quantity_presentation",
		"quantity_parsed",
		"quantity_parsed_left",
		"quantity_presentation_left",
		"lot_number",
		"catalog_number",
		"expires_at",
	).Create(&stock).Error
	if err != nil {
		return uuid.Nil, err
	}
	return stock.Id, nil
}

func (dataSrc StorehouseStocksDataSource) GetStorehouseStocksByUtiltyKey(key string) ([]models.StorehouseStock, error) {
	storehouseStocks := make([]models.StorehouseStock, 0)
	err := dataSrc.DbPsql.
		Where("storehouse_utility_key = ?", key).
		Find(&storehouseStocks).Error
	if err != nil {
		return storehouseStocks, err
	}
	return storehouseStocks, nil
}

func (dataSrc StorehouseStocksDataSource) GetStorehouseStockById(id string) (models.StorehouseStock, error) {
	var stock models.StorehouseStock
	err := dataSrc.DbPsql.
		Where("id = ?", id).
		First(&stock).Error
	if err != nil {
		return stock, err
	}
	return stock, nil
}

func (dataSrc StorehouseStocksDataSource) GetStorehouseInventoryStocks(filters models.StorehouseUtilitiesFilters) ([]models.StorehouseUtilityStockDetailed, int, error) {
	var totalRecords int64
	keysMatched := make([]string, 0)
	storehouseStocks := make([]models.StorehouseUtilityStockDetailed, 0)
	totalPages := 1
	sqlInstance := dataSrc.DbPsql.Table("storehouse_stocks")
	conditionStringFromJson := common.StructJsonSerializer(models.StorehouseUtilitiesFilters{
		UtilityKey:  filters.UtilityKey,
		UtilityName: filters.UtilityName,
	}, "json")
	if conditionStringFromJson != "" {
		// Get the utilities with key that matches the condition
		dataSrc.DbPsql.Select("key").Table("storehouse_utilities").Where(conditionStringFromJson).Find(&keysMatched)
		sqlInstance = sqlInstance.Where("storehouse_utility_key in (?)", keysMatched)
	}
	sqlInstance = sqlInstance.
		Limit(filters.Limit).
		Offset((filters.Page - 1) * filters.Limit).
		Preload("StorehouseUtility").
		Find(&storehouseStocks).
		Count(&totalRecords)

	if sqlInstance.Error != nil {
		return storehouseStocks, 0, sqlInstance.Error
	}

	totalPages = int(totalRecords) / filters.Limit
	if int(totalRecords)%filters.Limit != 0 {
		totalPages = totalPages + 1
	}

	return storehouseStocks, totalPages, nil
}

func (dataSrc StorehouseStocksDataSource) GetStorehouseInventoryUtilitiesDetailed(filters models.StorehouseUtilitiesFilters) ([]models.StorehouseUtilitsDetailedNoStocks, int, error) {
	var totalRecords int64
	keysMatched := make([]string, 0)
	totalPages := 1
	conditionStringFromJson := common.StructJsonSerializer(models.StorehouseUtilitiesFilters{
		UtilityKey:  filters.UtilityKey,
		UtilityName: filters.UtilityName,
	}, "json")
	// Get the utilities with key that matches the condition
	dataSrc.DbPsql.Select("key").Table("storehouse_utilities").Where(conditionStringFromJson).Find(&keysMatched)
	storehouseUtilitiesNoStocks := make([]models.StorehouseUtilitsDetailedNoStocks, 0)
	// Subquery to get total pieces left of a medicine
	subQuerySelect := dataSrc.DbPsql.
		Select("SUM(quantity_presentation) - SUM(quantity_presentation_used) as total_quantity_presentation_left").
		Where("storehouse_utility_key = ss.storehouse_utility_key").
		Table("storehouse_stocks")
	// Query medicines that match with the condition
	querySelect := dataSrc.DbPsql.
		Select("storehouse_utility_key,(?) as total_quantity_presentation_left", subQuerySelect).
		Group("ss.storehouse_utility_key").
		Table("storehouse_stocks as ss").
		Where("storehouse_utility_key in (?)", keysMatched)
	res := dataSrc.DbPsql.
		Omit("query_result.total_quantity_parsed_left").
		Table("(?) as query_result", querySelect).
		Preload("StorehouseUtility").
		Where("query_result.total_quantity_presentation_left < 1000").
		Find(&storehouseUtilitiesNoStocks).
		Count(&totalRecords)

	totalPages = int(totalRecords) / filters.Limit
	if totalPages*filters.Limit != int(totalRecords) {
		totalPages += 1
	}

	if res.Error != nil {
		return storehouseUtilitiesNoStocks, 0, res.Error
	}
	return storehouseUtilitiesNoStocks, totalPages, nil
}

func (dataSrc StorehouseStocksDataSource) UpdateStorehouseStock(id string, stock models.StorehouseStock) error {
	err := dataSrc.DbPsql.
		Select(
			"quantity_parsed",
			"quantity_presentation",
			"updated_at",
			"storehouse_utility_key",
			"lot_number",
			"catalog_number",
			"expires_at",
		).
		Where("id = ?", id).
		Updates(&stock).Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc StorehouseStocksDataSource) DeleteStorehouseStock(id string) error {
	err := dataSrc.DbPsql.
		Where("id = ?", id).
		Delete(&models.StorehouseStock{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc StorehouseStocksDataSource) IsStockUsed(id string) bool {
	var stock models.StorehouseStock
	dataSrc.DbPsql.
		Where("id = ?", id).
		Find(&stock)
	return stock.QuantityParsedUsed != 0 && stock.QuantityPresentationUsed != 0
}

func (dataSrc StorehouseStocksDataSource) IsStorehouseStockWithLotNumber(lot_number string, id string) bool {
	var storehouseStock models.StorehouseStock
	whereValues := []interface{}{fmt.Sprintf("%v", lot_number)}
	whereCondition := "\"lot_number\" = ?"
	if id != "" {
		whereCondition += " AND \"id\" != ?"
		whereValues = append(whereValues, fmt.Sprintf("%v", id))
	}
	err := dataSrc.DbPsql.
		Where(whereCondition, whereValues...).
		Take(&storehouseStock).
		Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (dataSrc StorehouseStocksDataSource) IsStorehouseStockWithCatalogNumber(catalog_number string, id string) bool {
	var storehouseStock models.StorehouseStock
	whereValues := []interface{}{fmt.Sprintf("%v", catalog_number)}
	whereCondition := "\"catalog_number\" = ?"
	if id != "" {
		whereCondition += " AND \"id\" != ?"
		whereValues = append(whereValues, fmt.Sprintf("%v", id))
	}
	err := dataSrc.DbPsql.
		Where(whereCondition, whereValues...).
		Take(&storehouseStock).
		Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
