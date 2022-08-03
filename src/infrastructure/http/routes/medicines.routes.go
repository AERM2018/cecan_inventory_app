package routes

import (
	"cecan_inventory/src/adapters/controllers"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"
	"cecan_inventory/src/infrastructure/http/middlewares"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitMedicinesRoutes(router router.Party, dbPsql *gorm.DB) {
	medicines := router.Party("/medicines")
	medicinesDataSource := datasources.MedicinesDataSource{DbPsql: dbPsql}
	val := middlewares.DbValidator{MedicineDataSrc: medicinesDataSource}
	controller := controllers.MedicinesController{}
	controller.New(medicinesDataSource)

	medicines.Get("/", controller.GetMedicinesCatalog)
	medicines.Post("/", val.IsMedicineInCatalogByKey, val.IsMedicineWithName, controller.InsertMedicineIntoCatalog)
	medicines.Put("/{key:string}", val.IsMedicineInCatalogByKey, val.IsMedicineWithName, controller.UpdateMedicine)
	medicines.Put("/{key:string}/reactivate", val.IsMedicineInCatalogByKey, controller.ReactivateMedicine)
	medicines.Delete("/{key:string}", val.IsMedicineInCatalogByKey, controller.DeleteMedicine)
}
