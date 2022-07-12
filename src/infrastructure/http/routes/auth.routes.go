package routes

import (
	"cecan_inventory/src/adapters/controllers"

	"github.com/kataras/iris/v12/core/router"
)

func InitUsersRoutes(router router.Party) {
	auth := router.Party("/auth")
	controller := controllers.AuthController{}
	auth.Post("/login", controller.Login)
}
