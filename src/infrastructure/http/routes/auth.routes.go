package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"
	customreqvalidations "cecan_inventory/infrastructure/http/middlewares/customReqValidations"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitUsersRoutes(router router.Party, dbPsql *gorm.DB) {
	auth := router.Party("/auth")
	userDataSource := datasources.UserDataSource{DbPsql: dbPsql}
	roleDataSource := datasources.RolesDataSource{DbPsql: dbPsql}
	controller := controllers.AuthController{}
	controller.New(userDataSource)
	val := middlewares.DbValidator{RolesDataSource: roleDataSource, UserDataSource: userDataSource}
	// Enpoints definition by HTTP method
	// Apply custom validations to the requests' body
	auth.Post("/login", controller.Login)

	auth.Post("/signup",
		middlewares.ValidateRequest(customreqvalidations.ValidateUser),
		val.IsRoleId, val.IsEmail, controller.SignUp)

	auth.Post("/renew", middlewares.VerifyJWT, controller.RenewToken)
}
