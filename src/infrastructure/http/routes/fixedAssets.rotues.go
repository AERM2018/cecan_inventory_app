package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitFixedAssetsRoutes(router router.Party, dbPsql *gorm.DB) {
	fixedAssets := router.Party("/fixed_assets")
	fixedAssetsDataSource := datasources.FixedAssetsDataSource{DbPsql: dbPsql}
	fixedAssetDescriptionsDataSource := datasources.FixedAssetDescriptionDataSource{DbPsql: dbPsql}
	controller := controllers.FixedAssetsController{
		FixedAssetsDataSource:            fixedAssetsDataSource,
		FixedAssetDescriptionsDataSource: fixedAssetDescriptionsDataSource,
	}
	controller.New()
	fixedAssets.Get("/", controller.GetFixedAssets)
	fixedAssets.Post("/", controller.CreateFixedAsset)
}
