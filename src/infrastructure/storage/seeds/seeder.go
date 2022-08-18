package seeds

import (
	"gorm.io/gorm"
)

// Seed type
type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

func All() []Seed{
	return []Seed{
		{Name:"CreateAdminRole",Run:func(db *gorm.DB) error {return CreateRoles(db,"Admin")}},
	}
}