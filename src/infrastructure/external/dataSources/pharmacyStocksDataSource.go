package datasources

import (
	"cecan_inventory/domain/models"

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
	res := dataSrc.DbPsql.First(&pharmacyStock, "id = ?", id)
	if res.Error != nil {
		return pharmacyStock, res.Error
	}
	return pharmacyStock, nil
}

func (dataSrc PharmacyStocksDataSource) GetPharmacyStocksByMedicineKey(medicineKey string) ([]models.PharmacyStock, error) {
	var pharmacyStocks []models.PharmacyStock
	res := dataSrc.DbPsql.Where("medicine_key = ?", medicineKey).Omit("created_at", "updated_at", "deleted_at").Find(&pharmacyStocks)
	if res.Error != nil {
		return pharmacyStocks, res.Error
	}
	return pharmacyStocks, nil
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
	res := dataSrc.DbPsql.First(&pharmacyStock)
	if res.Error != nil {
		return false, res.Error
	}
	isStockUsed := pharmacyStock.Pieces_used > 0 && pharmacyStock.CreatedAt.Unix() != pharmacyStock.UpdatedAt.Unix()
	return isStockUsed, nil
}
