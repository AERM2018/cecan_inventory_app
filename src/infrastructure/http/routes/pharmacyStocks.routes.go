package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitPharmacyStocksRoutes(router router.Party, dbPsql *gorm.DB) {
	pharmacyInventory := router.Party("/pharmacy_inventory")
	pharmacyStocksDataSource := datasources.PharmacyStocksDataSource{DbPsql: dbPsql}
	medicinesDataSource := datasources.MedicinesDataSource{DbPsql: dbPsql}
	// val := middlewares.DbValidator{PharmacyDataSrc: pharmacyStocksDataSource, MedicineDataSrc: medicinesDataSource}
	controller := controllers.PharmacyStocksController{
		PharmacyStocksDataSource: pharmacyStocksDataSource,
		MedicineDataSource:       medicinesDataSource,
	}
	controller.New()
	pharmacyInventory.Use(middlewares.VerifyJWT)
	pharmacyInventory.Get("/", controller.GetPharmacyStocks)
	pharmacyInventory.Put("/{id:string}", controller.UpdatePharmacyStock)

}
