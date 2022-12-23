package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitUserRuotes(router router.Party, dbPsql *gorm.DB) {
	users := router.Party("/users")
	usersDataSource := datasources.UserDataSource{DbPsql: dbPsql}
	usersController := controllers.UsersController{
		UserDataSource: usersDataSource,
	}
	usersController.New()
	// Declare dbvalidator and pass the correspond data source
	// val := middlewares.DbValidator{
	// 	StorehouseUtilityDataSource: usersDataSource,
	// }
	users.Use(middlewares.VerifyJWT)
	// Enpoints definition by HTTP methods
	// GET
	users.Get("/", usersController.GetUsers)
}
