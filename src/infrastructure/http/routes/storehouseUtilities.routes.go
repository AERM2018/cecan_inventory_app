package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"
	customreqvalidations "cecan_inventory/infrastructure/http/middlewares/customReqValidations"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitStorehouseUtilitiesRoutes(router router.Party, dbPsql *gorm.DB) {
	storehouseUtilities := router.Party("/storehouse_utilities")
	storehouseUtilitiesDataSource := datasources.StorehouseUtilitiesDataSource{DbPsql: dbPsql}
	storehouseUtilitesController := controllers.StorehouseUtilitiesController{StorehouseUtilitiesDataSource: storehouseUtilitiesDataSource}
	storehouseUtilitesController.New()
	// Declare dbvalidator and pass the correspond data source
	val := middlewares.DbValidator{
		StorehouseUtilityDataSource: storehouseUtilitiesDataSource,
	}
	storehouseUtilities.Use(middlewares.VerifyJWT)
	// Enpoints definition by HTTP methods
	storehouseUtilities.Get("/categories", storehouseUtilitesController.GetStorehouseUtilityCategories)
	storehouseUtilities.Get("/presentations", storehouseUtilitesController.GetStorehouseUtilityPresentations)
	storehouseUtilities.Get("/units", storehouseUtilitesController.GetStorehouseUtilityUnits)
	storehouseUtilities.Post("/",
		val.CanUserDoAction("Almacen"),
		val.IsStorehouseUtilityWithKey,
		middlewares.ValidateRequest(customreqvalidations.ValidateStorehouseUtility),
		storehouseUtilitesController.CreateStorehouseUtility,
	)
}
