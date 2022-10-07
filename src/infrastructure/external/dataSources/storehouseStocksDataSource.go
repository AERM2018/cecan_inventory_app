package datasources

import (
	"cecan_inventory/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorehouseStocksDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc StorehouseStocksDataSource) CreateStorehouseStock(stock models.StorehouseStock) (uuid.UUID, error) {
	err := dataSrc.DbPsql.Select("storehouse_utility_key", "quantity_presentation", "quantity_parsed").Create(&stock).Error
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

func (dataSrc StorehouseStocksDataSource) UpdateStorehouseStock(id string, stock models.StorehouseStock) error {
	err := dataSrc.DbPsql.
		Select("quantity_parsed", "quantity_presentation", "updated_at", "storehouse_utility_key").
		Where("id = ?", id).
		Updates(&stock).Error
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
