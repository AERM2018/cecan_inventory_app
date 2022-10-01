package seeds

import (
	"cecan_inventory/domain/mocks"
	"os"

	"gorm.io/gorm"
)

func CreateUsers(db *gorm.DB) error {
	if os.Getenv("GO_ENV") == "TEST" {
		roles := mocks.GetRolesMock("")
		for _, role := range roles {
			user := mocks.GetUserMock(role.Id.String(), 10)
			err := db.Create(&user).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
