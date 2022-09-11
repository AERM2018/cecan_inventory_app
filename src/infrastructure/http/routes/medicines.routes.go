package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitMedicinesRoutes(router router.Party, dbPsql *gorm.DB) {
	medicines := router.Party("/medicines")
	medicinesDataSource := datasources.MedicinesDataSource{DbPsql: dbPsql}
	pharmacyStocksDataSource := datasources.PharmacyStocksDataSource{DbPsql: dbPsql}
	val := middlewares.DbValidator{MedicineDataSrc: medicinesDataSource, PharmacyDataSrc: pharmacyStocksDataSource}
	controller := controllers.MedicinesController{MedicinesDataSource: medicinesDataSource, PharmacyStocksDataSource: pharmacyStocksDataSource}
	controller.New()

	medicines.Get("/", controller.GetMedicinesCatalog)
	medicines.Post("/", val.IsMedicineInCatalogByKey, val.IsMedicineWithName, controller.InsertMedicineIntoCatalog)
	medicines.Post("/{key:string}/pharmacy_inventory", val.IsMedicineInCatalogByKey, controller.InsertPharmacyStockOfMedicine)
	medicines.Put("/{key:string}", val.IsMedicineInCatalogByKey, val.IsMedicineWithName, val.IsMedicineWithKey, controller.UpdateMedicine)
	medicines.Put("/{key:string}/reactivate", val.IsMedicineInCatalogByKey, val.IsMedicineDeleted, controller.ReactivateMedicine)
	medicines.Delete("/{key:string}", val.IsMedicineInCatalogByKey, controller.DeleteMedicine)
}
