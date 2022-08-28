package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitPharmacyStocksRoutes(router router.Party, dbPsql *gorm.DB) {
	pharmacy_inventory := router.Party("/pharmacy_inventory")
	pharmacyStocksDataSource := datasources.PharmacyStocksDataSource{DbPsql: dbPsql}
	medicinesDataSource := datasources.MedicinesDataSource{DbPsql: dbPsql}
	// val := middlewares.DbValidator{PharmacyDataSrc: PharmacyStocksDataSource}
	controller := controllers.PharmacyStocksController{}
	controller.New(pharmacyStocksDataSource, medicinesDataSource)

	pharmacy_inventory.Get("/", controller.GetPharmacyStocks)
	pharmacy_inventory.Post("/", controller.InsertStockOfMedicine)
	pharmacy_inventory.Put("/{id:string}", controller.UpdatePharmacyStock)

}
