package routes

import (
	"cecan_inventory/src/adapters/controllers"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitMedicinesRoutes(router router.Party, dbPsql *gorm.DB) {
	medicines := router.Party("/medicines")
	medicinesDataSource := datasources.MedicinesDataSource{DbPsql: dbPsql}
	controller := controllers.MedicinesController{}
	controller.New(medicinesDataSource)

	medicines.Get("/", controller.GetMedicinesCatalog)
	medicines.Post("/", controller.InsertMedicineIntoCatalog)
}
