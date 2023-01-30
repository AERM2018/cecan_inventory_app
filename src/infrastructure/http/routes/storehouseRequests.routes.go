package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"
	customreqvalidations "cecan_inventory/infrastructure/http/middlewares/customReqValidations"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitStorehouseRequestsRoutes(router router.Party, dbPsql *gorm.DB) {
	storehouseRequests := router.Party("/storehouse/requests")
	storehouseRequestsDataSource := datasources.StorehouseRequestsDataSource{DbPsql: dbPsql}
	storehouseUtilitiesDataSource := datasources.StorehouseUtilitiesDataSource{DbPsql: dbPsql}
	storehouseRequestsController := controllers.StorehouseRequestsController{
		StorehouseRequestsDataSource: storehouseRequestsDataSource,
	}
	val := middlewares.DbValidator{
		StorehouseUtilityDataSource:  storehouseUtilitiesDataSource,
		StorehouseRequestsDataSource: storehouseRequestsDataSource,
	}
	storehouseRequestsController.New()
	storehouseRequests.Use(middlewares.VerifyJWT)
	storehouseRequests.Get("/", storehouseRequestsController.GetStorehouseRequests)
	storehouseRequests.Get("/{id:string}", val.IsStorehouseRequest, storehouseRequestsController.GetStorehouseRequestById)

	storehouseRequests.Post("/",
		middlewares.ValidateRequest(customreqvalidations.ValidateStorehouseRequest, "pieces_supplied"),
		val.AreStorehouseRequestItemsValid,
		storehouseRequestsController.CreateStorehouseRequest)

	storehouseRequests.Put("/{id:string}",
		middlewares.ValidateRequest(customreqvalidations.ValidateStorehouseRequest, "pieces_supplied"),
		val.IsStorehouseRequest,
		val.IsSameRequestCreator,
		val.IsRequestDeterminedStatus("Pendiente"),
		val.AreStorehouseRequestItemsValid,
		storehouseRequestsController.UpdateStorehouseRequest)

	storehouseRequests.Put("/{id:string}/complete",
		middlewares.ValidateRequest(customreqvalidations.ValidateStorehouseRequest, "pieces"),
		val.CanUserDoAction("Almacen"),
		val.IsStorehouseRequest,
		val.IsRequestDeterminedStatus("Pendiente", "Incompleta"),
		storehouseRequestsController.CompleteStorehouseRequest)

	storehouseRequests.Delete("/{id:string}",
		val.IsStorehouseRequest,
		val.IsSameRequestCreator,
		val.IsRequestDeterminedStatus("Pendiente"),
		storehouseRequestsController.DeleteStorehouseRequest)
}
