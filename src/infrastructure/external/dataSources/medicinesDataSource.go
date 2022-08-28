package datasources

import (
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

type MedicinesDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc MedicinesDataSource) InsertMedicine(medicine models.Medicine) error {
	res := dataSrc.DbPsql.Create(&medicine)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (dataSrc MedicinesDataSource) GetMedicineByKey(key string) (models.Medicine, error) {
	var medicineFound models.Medicine
	res := dataSrc.DbPsql.Find(&medicineFound, key)
	if res.Error != nil {
		return medicineFound, res.Error
	}
	return medicineFound, nil
}

func (dataSrc MedicinesDataSource) GetMedicinesCatalog() ([]models.Medicine, error) {
	var medicinesCatalog []models.Medicine
	res := dataSrc.DbPsql.Omit("created_at", "updated_at", "deletet_at").Find(&medicinesCatalog)
	if res.Error != nil {
		return medicinesCatalog, res.Error
	}
	return medicinesCatalog, nil
}

func (dataSrc MedicinesDataSource) UpdateMedicine(key string, medicine models.Medicine) (string, error) {
	res := dataSrc.DbPsql.Model(&models.Medicine{Key: key}).Updates(medicine)
	if res.Error != nil {
		return "", res.Error
	}
	return medicine.Key, nil
}

func (dataSrc MedicinesDataSource) DeleteMedicineByKey(key string) error {
	res := dataSrc.DbPsql.Where("key = ?", key).Delete(&models.Medicine{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (dataSrc MedicinesDataSource) ReactivateMedicine(key string) error {
	res := dataSrc.DbPsql.Unscoped().Model(&models.Medicine{}).Where("key = ?", key).Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
