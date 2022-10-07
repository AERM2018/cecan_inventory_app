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
