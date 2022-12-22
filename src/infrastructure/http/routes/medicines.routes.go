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
	// medicines.Use(val.CanUserDoAction("Farmacia"))
	// Enpoints definition by HTTP method
	// Apply custom validations to the requests' body
	medicines.Get("/", val.CanUserDoAction("Farmacia", "Medico"), controller.GetMedicinesCatalog)

	medicines.Post("/",
		middlewares.ValidateRequest(customreqvalidations.ValidateMedicine),
		val.CanUserDoAction("Farmacia"), val.IsMedicineWithKey, val.IsMedicineWithName, controller.InsertMedicineIntoCatalog)

	medicines.Post("/{key:string}/pharmacy_inventory",
		middlewares.ValidateRequest(customreqvalidations.ValidatePharmacyStock, "medicine_key"),
		val.CanUserDoAction("Farmacia"), val.IsMedicineInCatalogByKey, val.IsMedicineDeleted, val.IsPharmacyStockWithLotNumber, controller.InsertPharmacyStockOfMedicine)

	medicines.Put("/{key:string}",
		middlewares.ValidateRequest(customreqvalidations.ValidateMedicine),
		val.CanUserDoAction("Farmacia"), val.IsMedicineInCatalogByKey, val.IsMedicineWithName, val.IsMedicineWithKey, controller.UpdateMedicine)

	medicines.Put("/{key:string}/reactivate", val.CanUserDoAction("Farmacia"), val.IsMedicineInCatalogByKey, val.IsMedicineDeleted, controller.ReactivateMedicine)

	medicines.Delete("/{key:string}", val.CanUserDoAction("Farmacia"), val.IsMedicineInCatalogByKey, controller.DeleteMedicine)
}
