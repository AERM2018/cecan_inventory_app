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

func (dataSrc PharmacyStocksDataSource) GetPharmacyStocks(filters models.MedicinesFilters) ([]models.PharmacyStocksDetails, int,error) {
	var totalRecords int64
	pharmacyStocks := make([]models.PharmacyStocksDetails, 0)
	totalPages := 0
	// res := dataSrc.DbPsql.Raw(fmt.Sprintf("SELECT * FROM public.get_pharmacy_stocks_sorted_no_color('%v');", medicineKey)).Scan(&pharmacyStocks)
	res := dataSrc.DbPsql.Table("pharmacy_stocks")

	if filters.MedicineKey != "" {
		res = res.Where(fmt.Sprintf("medicine_key LIKE %v%v%v","'%",filters.MedicineKey,"%'"))
	}

	res.Count(&totalRecords)
	res = res.
		Limit(filters.Limit).
		Offset((filters.Page-1) * filters.Limit).
		Preload("Medicine").
		Find(&pharmacyStocks)

	if res.Error != nil {
		return pharmacyStocks,0, res.Error
	}

	totalPages = int(totalRecords) / filters.Limit
	if int(totalRecords) % filters.Limit != 0 {
		totalPages = totalPages + 1
	}
	
	return pharmacyStocks, totalPages, nil
}

func (dataSrc PharmacyStocksDataSource) GetMedicinesWithLessStockQty(filters models.MedicinesFilters)([]models.PharmacyStock, error){
	conditionString := ""
	conditionStringFromJson := common.StructJsonSerializer(models.MedicinesFilters{
		MedicineKey: filters.MedicineKey,
	})
	if conditionStringFromJson != ""{
		conditionString += "WHERE " + conditionStringFromJson
	}
	sqlQuery := fmt.Sprintf("SELECT medicine_key, (SELECT SUM(pieces) FROM pharmacy_stocks where medicine_key = ps.medicine_key) as pieces_left FROM pharmacy_stocks ps %v GROUP by (medicine_key) LIMIT %v OFFSET %v;", conditionString,filters.Limit, (filters.Page - 1) * filters.Limit)
	medicinesWithStocks := make([]models.PharmacyStock,0)
	res := dataSrc.DbPsql.Raw(sqlQuery).Joins("Medicine").Scan(&medicinesWithStocks)
	if res.Error != nil {
		return medicinesWithStocks, res.Error
	}
	return medicinesWithStocks, nil
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
