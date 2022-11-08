package config

import (
	"fmt"
	"os"

	"cecan_inventory/infrastructure/http/routes"
	"cecan_inventory/infrastructure/storage"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
)

type Server struct {
	IrisApp *iris.Application
	DbPsql  *gorm.DB
	Router  router.Party
	Port    string
}

func (server *Server) New() *iris.Application {
	server.IrisApp = iris.New()
	server.ConnectDatabase()
	server.setUpMiddlewares()
	server.Router = server.IrisApp.Party("/api/v1")
	routes.InitUsersRoutes(server.Router, server.DbPsql)
	routes.InitMedicinesRoutes(server.Router, server.DbPsql)
	routes.InitPharmacyStocksRoutes(server.Router, server.DbPsql)
	routes.InitPrescriptionsRoutes(server.Router, server.DbPsql)
	routes.InitRolesRoutes(server.Router, server.DbPsql)
	routes.InitStorehouseUtilitiesRoutes(server.Router, server.DbPsql)
	routes.InitStorehouseStocksRoutes(server.Router, server.DbPsql)
	routes.InitStorehouseRequestsRoutes(server.Router, server.DbPsql)
	routes.InitDepartmentsRoutes(server.Router, server.DbPsql)
	// Set port
	server.Port = os.Getenv("PORT")

	return server.IrisApp
}

func (server *Server) ConnectDatabase() {
	var errPsql error
	server.DbPsql, errPsql = storage.Connect()
	if errPsql != nil {
		panic(errPsql)
	}
	fmt.Println("PSQL online")
}

func (server *Server) setUpMiddlewares() {
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	server.IrisApp.Use(crs)
	server.IrisApp.AllowMethods(iris.MethodOptions)
}
func (server *Server) SetUpRouter() {
	server.Router = server.IrisApp.Party("/api/v1")
}

func (server *Server) StartListening() {
	server.IrisApp.Run(iris.Addr(":"+server.Port), iris.WithoutBodyConsumptionOnUnmarshal)
}
