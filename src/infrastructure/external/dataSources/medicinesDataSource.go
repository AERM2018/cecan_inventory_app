package datasources

import (
	"cecan_inventory/domain/common"
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

func (dataSrc MedicinesDataSource) GetMedicinesCatalog(filters models.MedicinesFilters, includeDeleted bool) ([]models.Medicine, int, error) {
	var totalRecords int64
	medicinesCatalog := make([]models.Medicine, 0)
	conditionStringFromJson := common.StructJsonSerializer(models.MedicinesFilters{
		MedicineKey:  filters.MedicineKey,
		MedicineName: filters.MedicineName,
	}, "json", "OR")
	res := dataSrc.DbPsql.Table("medicines")
	if includeDeleted {
		res = res.Unscoped()
	}

	res = res.
		Where(conditionStringFromJson).
		Count(&totalRecords).
		Limit(filters.Limit).
		Offset((filters.Page - 1) * filters.Limit).
		Find(&medicinesCatalog)

	totalPages := int(totalRecords) / filters.Limit
	if totalPages*filters.Limit != int(totalRecords) {
		totalPages += 1
	}
	return medicinesCatalog, totalPages, nil
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
