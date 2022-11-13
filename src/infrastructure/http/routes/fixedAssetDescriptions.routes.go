package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitFixedAssetDescriptionsRoutes(router router.Party, dbPsql *gorm.DB) {
	fixedAssetDescriptions := router.Party("/fixed_asset_descriptions")
	fixedAssetDescriptionsDataSource := datasources.FixedAssetDescriptionDataSource{DbPsql: dbPsql}
	controller := controllers.FixedAssetDescriptionsController{FixedAssetDescriptionsDataSource: fixedAssetDescriptionsDataSource}
	controller.New()
	fixedAssetDescriptions.Get("/", controller.GetFixedAssetDescriptions)
}
