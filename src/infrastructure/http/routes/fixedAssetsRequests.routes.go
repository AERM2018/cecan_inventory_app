package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitFixedAssetsRequestsRoutes(router router.Party, dbPsql *gorm.DB) {
	fixedAssetsRequests := router.Party("/fixed_assets_requests")
	// Datasources needed
	fixedAssetsRequestDataSource := datasources.FixedAssetsRequetsDataSource{DbPsql: dbPsql}
	fixedAssetsDataSource := datasources.FixedAssetsDataSource{DbPsql: dbPsql}
	fixedAssetsDescriptionsDataSource := datasources.FixedAssetDescriptionDataSource{DbPsql: dbPsql}
	// Db validator instance
	val := middlewares.DbValidator{
		FixedAssetsDataSource:         fixedAssetsDataSource,
		FixedAssetsRequestsDataSource: fixedAssetsRequestDataSource,
	}
	// Controller definition
	controller := controllers.FixedAssetsRequestsController{
		FixedAssetsRequestsDataSource:     fixedAssetsRequestDataSource,
		FixedAssetsDataSource:             fixedAssetsDataSource,
		FixedAssetsDescriptionsDataSource: fixedAssetsDescriptionsDataSource,
	}
	controller.New()
	fixedAssetsRequests.Use(middlewares.VerifyJWT)
	fixedAssetsRequests.Get("/", controller.GetFixedAssetsRequests)
	fixedAssetsRequests.Get("/{id:string}", val.FindFixedAssetsRequestById, controller.GetFixedAssetsRequestById)
	fixedAssetsRequests.Post("/departments/{departmentId:string}",
		val.AreFixedAssetsValidFromRequest,
		controller.CreateFixedAssetsRequest)
	fixedAssetsRequests.Delete("/{id:string}",
		val.FindFixedAssetsRequestById,
		controller.DeleteFixedAssetsRequest)
}
