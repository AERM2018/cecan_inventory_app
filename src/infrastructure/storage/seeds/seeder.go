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
		{Name: "CreateAdminRole", Run: func(db *gorm.DB) error { return CreateRoles(db, "Admin") }},
		{Name: "CreateAdminUser", Run: func(db *gorm.DB) error { return CreateAdminUser(db) }},
		{Name: "CreateMedicines", Run: func(db *gorm.DB) error { return CreateMedicines(db) }},
	}
}
