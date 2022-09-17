package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"
	customreqvalidations "cecan_inventory/infrastructure/http/middlewares/customReqValidations"

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
	// Use middlewares for all the routes
	medicines.Use(middlewares.VerifyJWT)
	medicines.Use(val.CanUserDoAction("Farmacia"))
	// Enpoints definition by HTTP method
	// Apply custom validations to the requests' body
	medicines.Get("/", controller.GetMedicinesCatalog)

	medicines.Post("/",
		middlewares.ValidateRequest(customreqvalidations.ValidateMedicine),
		val.IsMedicineWithKey, val.IsMedicineWithName, controller.InsertMedicineIntoCatalog)

	medicines.Post("/{key:string}/pharmacy_inventory",
		middlewares.ValidateRequest(customreqvalidations.ValidatePharmacyStock, "medicine_key"),
		val.IsMedicineInCatalogByKey, val.IsMedicineDeleted, controller.InsertPharmacyStockOfMedicine)

	medicines.Put("/{key:string}",
		middlewares.ValidateRequest(customreqvalidations.ValidateMedicine),
		val.IsMedicineInCatalogByKey, val.IsMedicineWithName, val.IsMedicineWithKey, controller.UpdateMedicine)

	medicines.Put("/{key:string}/reactivate", val.IsMedicineInCatalogByKey, val.IsMedicineDeleted, controller.ReactivateMedicine)

	medicines.Delete("/{key:string}", val.IsMedicineInCatalogByKey, controller.DeleteMedicine)
}
