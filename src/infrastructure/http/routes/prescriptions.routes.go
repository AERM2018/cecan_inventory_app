package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitPrescriptionsRoutes(router router.Party, dbPsql *gorm.DB) {
	prescriptions := router.Party("/prescriptions")
	prescriptionsDataSource := datasources.PrescriptionsDataSource{DbPsql: dbPsql}
	pharmacyStocksDataSource := datasources.PharmacyStocksDataSource{DbPsql: dbPsql}
	roleDataSource := datasources.RolesDataSource{DbPsql: dbPsql}
	val := middlewares.DbValidator{
		RolesDataSource:        roleDataSource,
		PrescriptionDataSource: prescriptionsDataSource,
		PharmacyDataSrc:        pharmacyStocksDataSource,
	}
	controller := controllers.PrescriptionsController{
		PrescriptionsDataSource: prescriptionsDataSource}
	controller.New()

	prescriptions.Use(middlewares.VerifyJWT)
	prescriptions.Get("/", val.CanUserDoAction("Medico"), controller.GetPrescriptions)
	prescriptions.Get("/{id:string}", val.CanUserDoAction("Medico", "Farmacia"), controller.GetPrescriptionById)
	prescriptions.Post("/", val.CanUserDoAction("Medico"), controller.CreatePrescription)
	prescriptions.Put("/{id:string}", val.CanUserDoAction("Medico"), controller.UpdatePrescription)
}
