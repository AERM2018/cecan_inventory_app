package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitPrescriptionsRoutes(router router.Party, dbPsql *gorm.DB) {
	prescriptions := router.Party("/prescriptions")
	prescriptionsDataSource := datasources.PrescriptionsDataSource{DbPsql: dbPsql}
	controller := controllers.PrescriptionsController{}
	controller.New(prescriptionsDataSource)

	prescriptions.Post("/", controller.CreatePrescription)
	prescriptions.Get("/{id:string}", controller.GetPrescriptionById)
}
