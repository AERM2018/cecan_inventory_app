package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitFixedAssetsRoutes(router router.Party, dbPsql *gorm.DB) {
	fixedAssets := router.Party("/fixed_assets")
	// Datasources
	fixedAssetsDataSource := datasources.FixedAssetsDataSource{DbPsql: dbPsql}
	fixedAssetDescriptionsDataSource := datasources.FixedAssetDescriptionDataSource{DbPsql: dbPsql}
	// Db validator instance
	val := middlewares.DbValidator{
		FixedAssetsDataSource: fixedAssetsDataSource,
	}
	// Controller's definition
	controller := controllers.FixedAssetsController{
		FixedAssetsDataSource:            fixedAssetsDataSource,
		FixedAssetDescriptionsDataSource: fixedAssetDescriptionsDataSource,
	}
	controller.New()
	fixedAssets.Get("/", controller.GetFixedAssets)
	fixedAssets.Post("/", controller.CreateFixedAsset)
	fixedAssets.Post("/file", controller.UploadFixedAssetsFromFile)
	fixedAssets.Put("/{key:string}",
		val.FindFixedAssetByKey,
		val.IsFixedAssetWithKey,
		val.IsFixedAssetWithSeries,
		controller.UpdateFixedAsset)
	fixedAssets.Delete("/{key:string}",
		val.FindFixedAssetByKey,
		controller.DeleteFixedAsset)

}
