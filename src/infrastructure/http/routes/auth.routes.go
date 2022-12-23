package routes

import (
	"cecan_inventory/adapters/controllers"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"cecan_inventory/infrastructure/http/middlewares"
	customreqvalidations "cecan_inventory/infrastructure/http/middlewares/customReqValidations"

	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

func InitAuthRoutes(router router.Party, dbPsql *gorm.DB) {
	auth := router.Party("/auth")
	userDataSource := datasources.UserDataSource{DbPsql: dbPsql}
	passwordResetCodesDataSource := datasources.PasswordResetCodesDataSource{DbPsql: dbPsql}
	roleDataSource := datasources.RolesDataSource{DbPsql: dbPsql}
	controller := controllers.AuthController{PasswordResetCodesDataSource: passwordResetCodesDataSource}
	controller.New(userDataSource)
	val := middlewares.DbValidator{RolesDataSource: roleDataSource, UserDataSource: userDataSource}
	// Enpoints definition by HTTP method
	// Apply custom validations to the requests' body
	auth.Post("/login", controller.Login)

	auth.Post("/signup",
		middlewares.ValidateRequest(customreqvalidations.ValidateUser),
		val.IsRoleId, val.IsEmail, controller.SignUp)

	auth.Post("/renew", middlewares.VerifyJWT, controller.RenewToken)
	auth.Post("/password_reset_code", val.IsEmail, controller.GeneratePasswordResetCode)
	auth.Post("/password_reset/users/{userId:string}", controller.ResetPassword)
}
