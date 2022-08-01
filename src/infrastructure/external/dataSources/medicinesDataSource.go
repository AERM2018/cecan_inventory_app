package datasources

import (
	"cecan_inventory/src/domain/models"

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
	res := dataSrc.DbPsql.Find(&medicinesCatalog)
	if res.Error != nil {
		return medicinesCatalog, res.Error
	}
	return medicinesCatalog, nil
}
