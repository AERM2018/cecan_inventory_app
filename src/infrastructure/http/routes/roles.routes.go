package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitRolesRoutes(router router.Party, dbPsql *gorm.DB) {
	roles := router.Party("/roles")
	rolesDataSource := datasources.RolesDataSource{DbPsql: dbPsql}
	controller := controllers.RolesController{RolesDataSource: rolesDataSource}
	controller.New()
	roles.Get("/", controller.GetRoles)
}
