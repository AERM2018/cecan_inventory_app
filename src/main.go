package main

import (
	"cecan_inventory/src/infrastructure/config"
	"cecan_inventory/src/infrastructure/http/routes"
)

func main() {
	server := config.Server{}
	server.New()
	server.ConnectDatabase()
	server.SetUpRouter()
	// Init all the routes
	routes.InitUsersRoutes(server.Router)
	// Start app
	server.StartListening()
}
