package main

import (
	"cecan_inventory/infrastructure/config"
)

func main() {
	server := config.Server{}
	server.New()
	server.StartListening()
}
