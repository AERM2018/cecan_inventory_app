package datasources

import (
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PharmacyStocksDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc PharmacyStocksDataSource) InsertStockOfMedicine(pharmacyStock models.PharmacyStock) (uuid.UUID, error) {
	res := dataSrc.DbPsql.Create(&pharmacyStock)
	if res.Error != nil {
		return uuid.Nil, res.Error
	}
	return pharmacyStock.Id, nil
}

func (dataSrc PharmacyStocksDataSource) GetPharmacyStockById(id uuid.UUID) (models.PharmacyStock, error) {
	var pharmacyStock models.PharmacyStock
	res := dataSrc.DbPsql.Preload("Medicine").First(&pharmacyStock, "id = ?", id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return pharmacyStock, res.Error
	}
	return pharmacyStock, nil
}

func (dataSrc PharmacyStocksDataSource) GetPharmacyStocks(filters models.MedicinesFilters) ([]models.PharmacyStocksDetails, int, error) {
	var totalRecords int64
	pharmacyStocks := make([]models.PharmacyStocksDetails, 0)
	keysMatched := make([]string, 0)
	totalPages := 0
	conditionStringFromJson := common.StructJsonSerializer(models.MedicinesFilters{
		MedicineKey:  filters.MedicineKey,
		MedicineName: filters.MedicineName,
	}, "json", "OR")
	dataSrc.DbPsql.Select("key").Table("medicines").Where(conditionStringFromJson).Find(&keysMatched)
	res := dataSrc.DbPsql.Table("pharmacy_stocks").Where("medicine_key in (?)", keysMatched)
	res = res.
		Limit(filters.Limit).
		Offset((filters.Page - 1) * filters.Limit).
		Preload("Medicine").
		Find(&pharmacyStocks).
		Count(&totalRecords)

	if res.Error != nil {
		return pharmacyStocks, 0, res.Error
	}

	totalPages = int(totalRecords) / filters.Limit
	if int(totalRecords)%filters.Limit != 0 {
		totalPages = totalPages + 1
	}

	return pharmacyStocks, totalPages, nil
}

func (dataSrc PharmacyStocksDataSource) GetMedicinesWithLessStockQty(filters models.MedicinesFilters) ([]models.PharmacyStocksDetailedNoStocks, int, error) {
	var totalRecords int64
	keysMatched := make([]string, 0)
	totalPages := 1
	conditionStringFromJson := common.StructJsonSerializer(models.MedicinesFilters{
		MedicineKey:  filters.MedicineKey,
		MedicineName: filters.MedicineName,
	}, "json", "OR")
	// Get the medicines keys that matches the condition
	dataSrc.DbPsql.Select("key").Table("medicines").Where(conditionStringFromJson).Find(&keysMatched)
	medicinesWithStocks := make([]models.PharmacyStocksDetailedNoStocks, 0)
	// Subquery to get total pieces left of a medicine
	subQuerySelect := dataSrc.DbPsql.
		Select("SUM(pieces) - SUM(pieces_used) as pieces_left").
		Where("medicine_key = ps.medicine_key").
		Table("pharmacy_stocks")
	// Query medicines that match with the condition
	querySelect := dataSrc.DbPsql.
		Select("medicine_key,(?) as total_pieces_left", subQuerySelect).
		Group("ps.medicine_key").
		Table("pharmacy_stocks as ps").
		Where("medicine_key in (?)", keysMatched)
	res := dataSrc.DbPsql.
		Omit("query_result.total_pieces").
		Table("(?) as query_result", querySelect).
		Joins("Medicine").
		Where("query_result.total_pieces_left < 1000").
		Find(&medicinesWithStocks).
		Count(&totalRecords)

	totalPages = int(totalRecords) / filters.Limit
	if totalPages*filters.Limit != int(totalRecords) {
		totalPages += 1
	}

	if res.Error != nil {
		return medicinesWithStocks, 0, res.Error
	}
	return medicinesWithStocks, totalPages, nil
}

func (dataSrc PharmacyStocksDataSource) UpdatePharmacyStock(id uuid.UUID, pharmacyStock models.PharmacyStockToUpdate) error {
	res := dataSrc.DbPsql.Table("pharmacy_stocks").Where("id = ?", id).Updates(&pharmacyStock)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (dataSrc PharmacyStocksDataSource) DeletePharmacyStock(id uuid.UUID, isPermanent bool) error {
	pointerDb := dataSrc.DbPsql
	if isPermanent {
		pointerDb = pointerDb.Unscoped()
	}
	res := pointerDb.Where("id = ?", id).Delete(&models.PharmacyStock{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (dataSrc PharmacyStocksDataSource) IsStockUsed(id uuid.UUID) (bool, error) {
	var pharmacyStock models.PharmacyStock
	res := dataSrc.DbPsql.Where("id = ?", id).First(&pharmacyStock)
	if res.Error != nil {
		return false, res.Error
	}
	isStockUsed := pharmacyStock.Pieces_used > 0
	return isStockUsed, nil
}

func (dataSrc PharmacyStocksDataSource) IsPharmacyStockWithLotNumber(lotNumber string, pharmacyStockId string) bool {
	var pharmacyStock models.PharmacyStock
	whereValues := []interface{}{fmt.Sprintf("'%v'", lotNumber)}
	whereCondition := "\"lot_number\" = ?"
	if pharmacyStockId != "" {
		whereCondition += " AND \"id\" != ?"
		whereValues = append(whereValues, fmt.Sprintf("%v", pharmacyStockId))
	}
	// fmt.Println(whereValues...)
	err := dataSrc.DbPsql.Where(whereCondition, whereValues...).Take(&pharmacyStock).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
