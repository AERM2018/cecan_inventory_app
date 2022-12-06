package seeds

import (
	"cecan_inventory/domain/mocks"

	"gorm.io/gorm"
)

func CreateResponsibleUsers(db *gorm.DB) error {
	users := mocks.GetDirectorAndSubDirectorUsersMockSeed()
	for _, user := range users {
		db.FirstOrCreate(&user, user)
	}
	return nil
}
