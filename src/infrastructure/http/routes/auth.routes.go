package routes

import (
	"cecan_inventory/src/adapters/controllers"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitUsersRoutes(router router.Party, dbPsql *gorm.DB) {
	auth := router.Party("/auth")
	userDataSource := datasources.UserDataSource{DbPsql:dbPsql}
	controller := controllers.AuthController{ UserDataSource:userDataSource}
	auth.Post("/login", controller.Login)
}
