package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"
	customreqvalidations "cecan_inventory/infrastructure/http/middlewares/customReqValidations"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitPharmacyStocksRoutes(router router.Party, dbPsql *gorm.DB) {
	pharmacyInventory := router.Party("/pharmacy_inventory")

	pharmacyStocksDataSource := datasources.PharmacyStocksDataSource{DbPsql: dbPsql}
	medicinesDataSource := datasources.MedicinesDataSource{DbPsql: dbPsql}
	val := middlewares.DbValidator{PharmacyDataSrc: pharmacyStocksDataSource, MedicineDataSrc: medicinesDataSource}
	controller := controllers.PharmacyStocksController{
		PharmacyStocksDataSource: pharmacyStocksDataSource,
		MedicineDataSource:       medicinesDataSource,
	}
	controller.New()
	// Use middlewares for all the routes
	pharmacyInventory.Use(middlewares.VerifyJWT)
	// Enpoints definition by HTTP method
	// Apply custom validations to the requests' body
	pharmacyInventory.Get("/", controller.GetPharmacyStocks)
	// pharmacyInventory.Get("/", websocket.Handler(websocketServer))

	pharmacyInventory.Put("/{id:string}",
		middlewares.ValidateRequest(customreqvalidations.ValidatePharmacyStock),
		val.CanUserDoAction("Farmacia"),
		val.IsPharmacyStockById,
		val.IsMedicineInCatalogByKey,
		val.IsPharmacyStockUsed,
		val.IsPharmacyStockWithLotNumber,
		controller.UpdatePharmacyStock)

	pharmacyInventory.Delete("/{id:string}", val.CanUserDoAction("Farmacia"), val.IsPharmacyStockById, val.IsPharmacyStockUsed, controller.DeletePharmacyStock)

}
