package seeds

import (
	"gorm.io/gorm"
)

// Seed type
type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

func All() []Seed {
	return []Seed{
		{Name: "CreateAdminRole", Run: func(db *gorm.DB) error { return CreateRoles(db) }},
		{Name: "CreateAdminUser", Run: func(db *gorm.DB) error { return CreateAdminUser(db) }},
		{Name: "CreateUsers", Run: func(db *gorm.DB) error { return CreateUsers(db) }},
		{Name: "CreatePrescriptionStatues", Run: func(db *gorm.DB) error { return CreatePrescriptionStatues(db) }},
		{Name: "CreateMedicines", Run: func(db *gorm.DB) error { return CreateMedicines(db) }},
		{Name: "CreatePharmacyStock", Run: func(db *gorm.DB) error { return CreatePharmacyStock(db) }},
		{Name: "CreatePrescriptions", Run: func(db *gorm.DB) error { return CreatePrescriptions(db) }},
	}
}
