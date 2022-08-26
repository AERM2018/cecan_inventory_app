package main

import (
	"cecan_inventory/infrastructure/config"
	"cecan_inventory/infrastructure/http/routes"
)

func main() {
	server := config.Server{}
	server.New()
	server.ConnectDatabase()
	server.SetUpRouter()
	// Init all the routes
	routes.InitUsersRoutes(server.Router, server.DbPsql)
	routes.InitMedicinesRoutes(server.Router, server.DbPsql)
	routes.InitPharmacyStocksRoutes(server.Router, server.DbPsql)
	routes.InitPrescriptionsRoutes(server.Router, server.DbPsql)
	// Start app
	server.StartListening()
}
