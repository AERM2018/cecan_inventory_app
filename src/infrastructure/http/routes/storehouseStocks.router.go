package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"
	customreqvalidations "cecan_inventory/infrastructure/http/middlewares/customReqValidations"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitStorehouseStocksRoutes(router router.Party, dbPsql *gorm.DB) {
	storehouseStocks := router.Party("/storehouse_inventory")
	storehouseStocksDataSource := datasources.StorehouseStocksDataSource{DbPsql: dbPsql}
	storehouseUtilitiesDataSource := datasources.StorehouseUtilitiesDataSource{DbPsql: dbPsql}
	storehouseStocksController := controllers.StorehouseStocksController{
		StorehouseStocksDataSource:    storehouseStocksDataSource,
		StorehouseUtilitiesDataSource: storehouseUtilitiesDataSource,
	}

	val := middlewares.DbValidator{
		StorehouseStocksDataSource:  storehouseStocksDataSource,
		StorehouseUtilityDataSource: storehouseUtilitiesDataSource,
	}

	storehouseStocksController.New()
	// Use middlewares for all the routes
	storehouseStocks.Use(middlewares.VerifyJWT)
	// Enpoints definition by HTTP methods
	// GET
	storehouseStocks.Get("/", storehouseStocksController.GetStorehouseInventory)
	// PUT
	storehouseStocks.Put("/{id:string}",
		val.CanUserDoAction("Almacen"),
		middlewares.ValidateRequest(customreqvalidations.ValidateStorehouseStock),
		val.FindStorehouseStockById,
		val.IsStorehouseStockUsed,
		val.FindStorehouseUtilityByKey,
		val.IsStorehouseStockWithLotNumber,
		val.IsStorehouseStockWithCatalogNumber,
		storehouseStocksController.UpdateStorehouseStock)
	// DELETE
	storehouseStocks.Delete("/{id:string}",
		val.CanUserDoAction("Almacen"),
		val.FindStorehouseStockById,
		val.IsStorehouseStockUsed,
		storehouseStocksController.DeleteStorehouseStock,
	)
}
