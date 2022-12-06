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
		{Name: "CreateSuperiorRoles", Run: func(db *gorm.DB) error { return CreateSuperiorRoles(db) }},
		{Name: "CreateAdminUser", Run: func(db *gorm.DB) error { return CreateAdminUser(db) }},
		{Name: "CreateUsers", Run: func(db *gorm.DB) error { return CreateUsers(db) }},
		{Name: "CreateDepartments", Run: func(db *gorm.DB) error { return CreateDepartments(db) }},
		{Name: "CreatePrescriptionStatues", Run: func(db *gorm.DB) error { return CreatePrescriptionStatues(db) }},
		{Name: "CreateMedicines", Run: func(db *gorm.DB) error { return CreateMedicines(db) }},
		{Name: "CreatePharmacyStock", Run: func(db *gorm.DB) error { return CreatePharmacyStock(db) }},
		{Name: "CreatePrescriptions", Run: func(db *gorm.DB) error { return CreatePrescriptions(db) }},
		{Name: "CreateStorehouseUtilityCategories", Run: func(db *gorm.DB) error { return CreateStorehouseUtilityCategories(db) }},
		{Name: "CreateStorehouseUtilityPresentations", Run: func(db *gorm.DB) error { return CreateStorehouseUtilityPresentations(db) }},
		{Name: "CreateStorehouseUtilityUnits", Run: func(db *gorm.DB) error { return CreateStorehouseUtilityUnits(db) }},
		{Name: "CreateStorehouseRequestStatuses", Run: func(db *gorm.DB) error { return CreateStorehouseRequestStatuses(db) }},
		{Name: "CreateFixedAssetDescriptions", Run: func(db *gorm.DB) error { return CreateFixedAssetDescriptions(db) }},
		{Name: "CreateDirectorAndSubDirectorUsers", Run: func(db *gorm.DB) error { return CreateResponsibleUsers(db) }},
	}
}
