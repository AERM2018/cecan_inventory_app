package datasources

import (
	"cecan_inventory/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrescriptionsDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc PrescriptionsDataSource) CreatePrescription(prescription models.Prescription) (uuid.UUID, error) {
	res := dataSrc.DbPsql.Create(&prescription)
	if res.Error != nil {
		return uuid.Nil, res.Error
	}
	return prescription.Id, nil
}

func (dataSrc PrescriptionsDataSource) GetPrescriptionById(id uuid.UUID) (models.PrescriptionDetialed, error) {
	var prescription models.PrescriptionDetialed
	res := dataSrc.DbPsql.Model(models.Prescription{}).Preload("PrescriptionStatus", func(db *gorm.DB) *gorm.DB {
		return db.Omit("created_at", "updated_at", "deleted_at")
	}).Preload("Medicines", func(db *gorm.DB) *gorm.DB {
		return db.Omit("id", "created_at", "updated_at", "deleted_at")
	}).Preload("Medicines.Medicine", func(db *gorm.DB) *gorm.DB {
		return db.Omit("created_at", "updated_at", "deleted_at")
	}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password", "email", "created_at", "updated_at", "deleted_at")
	}).Preload("User.Role").Where("id = ?", id).Take(&prescription)
	if res.Error != nil {
		return prescription, res.Error
	}
	return prescription, nil
}

func (dataSrc PrescriptionsDataSource) PutMedicineIntoPrescription(medicineForPrescription models.PrescriptionsMedicines, prescriptionId uuid.UUID) error {
	prescriptionsMedicines := models.PrescriptionsMedicines{MedicineKey: medicineForPrescription.MedicineKey, Pieces: medicineForPrescription.Pieces, PrescriptionId: prescriptionId}
	res := dataSrc.DbPsql.Create(&prescriptionsMedicines)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
