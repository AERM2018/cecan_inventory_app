package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitStorehouseRequestsRoutes(router router.Party, dbPsql *gorm.DB) {
	storehouseRequests := router.Party("/storehouse/requests")
	storehouseRequestsDataSource := datasources.StorehouseRequestsDataSource{DbPsql: dbPsql}
	storehouseRequestsController := controllers.StorehouseRequestsController{
		StorehouseRequestsDataSource: storehouseRequestsDataSource,
	}
	storehouseRequestsController.New()
	storehouseRequests.Use(middlewares.VerifyJWT)
	storehouseRequests.Get("/", storehouseRequestsController.GetStorehouseRequests)
	storehouseRequests.Get("/{id:string}", storehouseRequestsController.GetStorehouseRequestById)
	storehouseRequests.Post("/", storehouseRequestsController.CreateStorehouseRequest)
}
