package datasources

import (
	"cecan_inventory/domain/models"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrescriptionsDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc PrescriptionsDataSource) CreatePrescription(prescription models.Prescription, medicines []models.PrescriptionsMedicines) (uuid.UUID, error) {
	isErrInTransaction := dataSrc.DbPsql.Transaction(func(tx *gorm.DB) error {
		// Create prescription and get id
		createError := tx.Create(&prescription).Error
		if createError != nil {
			return errors.New("No se pudo crear una nueva receta, verifique los datos y vuelvalo a intentar.")
		}
		// Specify the amount of medicine needed for the prescription
		for _, medicineForPrescription := range medicines {
			prescriptionsMedicines := models.PrescriptionsMedicines{
				MedicineKey:    medicineForPrescription.MedicineKey,
				Pieces:         medicineForPrescription.Pieces,
				PrescriptionId: prescription.Id}
			prescriptionMedicineError := tx.Create(&prescriptionsMedicines).Error
			if prescriptionMedicineError != nil {
				return errors.New("No se pudo crear la receta debido a que no se pudo asignar los medicamentos a la misma.")
			}
		}
		// Reserve the medicine specified
		reservationError := tx.Exec(fmt.Sprintf("Call public.reserve_medicines_for_prescription('%v');", prescription.Id)).Error
		if reservationError != nil {
			var customError string
			_, customError, _ = strings.Cut(reservationError.Error(), "ERROR: ")
			customError, _, _ = strings.Cut(customError, " (SQLSTATE P0001)")
			customError += "."
			return errors.New(customError)

		}
		return nil
	})
	if isErrInTransaction != nil {
		return uuid.Nil, isErrInTransaction
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

func (dataSrc PrescriptionsDataSource) GetPrescriptions(userId string) ([]models.PrescriptionDetialed, error) {
	var prescriptions []models.PrescriptionDetialed
	dbQuery := dataSrc.DbPsql.Model(models.Prescription{}).Preload("PrescriptionStatus", func(db *gorm.DB) *gorm.DB {
		return db.Omit("created_at", "updated_at", "deleted_at")
	}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password", "email", "created_at", "updated_at", "deleted_at")
	}).Preload("User.Role")
	if userId != "" {
		dbQuery = dbQuery.Where("user_id = ?", userId)
	}
	dbQuery.Find(&prescriptions)
	if dbQuery.Error != nil {
		return prescriptions, dbQuery.Error
	}
	return prescriptions, nil
}

func (dataSrc PrescriptionsDataSource) PutMedicineIntoPrescription(medicineForPrescription models.PrescriptionsMedicines, prescriptionId uuid.UUID) error {
	prescriptionsMedicines := models.PrescriptionsMedicines{MedicineKey: medicineForPrescription.MedicineKey, Pieces: medicineForPrescription.Pieces, PrescriptionId: prescriptionId}
	res := dataSrc.DbPsql.Create(&prescriptionsMedicines)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
