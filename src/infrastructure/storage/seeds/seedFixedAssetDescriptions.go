package seeds

import (
	"cecan_inventory/domain/mocks"

	"gorm.io/gorm"
)

func CreateFixedAssetDescriptions(db *gorm.DB) error {
	fixedAssetDescriptions := mocks.GetFixedAssetsDescriptionsMockSeed()
	for _, fixAssetDescription := range fixedAssetDescriptions {
		db.FirstOrCreate(&fixAssetDescription, fixAssetDescription)
	}
	return nil
}
