package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitUsersRoutes(router router.Party, dbPsql *gorm.DB) {
	auth := router.Party("/auth")
	userDataSource := datasources.UserDataSource{DbPsql: dbPsql}
	roleDataSource := datasources.RolesDataSource{DbPsql: dbPsql}
	controller := controllers.AuthController{}
	controller.New(userDataSource)
	val := middlewares.DbValidator{RoleDataSource: roleDataSource, UserDataSource: userDataSource}
	auth.Post("/login", controller.Login)
	auth.Post("/signup", val.IsRoleId, val.IsEmail, controller.SignUp)
	auth.Post("/renew", middlewares.VerifyJWT, controller.RenewToken)
}
