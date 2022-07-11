package main

import (
	"cecan_inventory/src/infrastructure/config"
)

func main() {
	server := config.Server{};
	server.New();
	server.ConnectDatabase()
}