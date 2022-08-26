package config

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"cecan_inventory/infrastructure/storage"

	"github.com/joho/godotenv"
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

func (server *Server) New() {
	server.IrisApp = iris.New()
	_, filename, _, _ := runtime.Caller(0)
	if os.Getenv("GO_ENV") != "production" {
		// load env variables from a .env file
		envPath := path.Join(path.Dir(filename), "../../../.env")
		err := godotenv.Load(envPath)
		if err != nil {
			panic(err)
		}
	}
	// Set port
	server.Port = os.Getenv("PORT")
}

func (server *Server) ConnectDatabase() {
	var errPsql error
	server.DbPsql, errPsql = storage.Connect()
	if errPsql != nil {
		fmt.Println(errPsql)
	}
	fmt.Println("PSQL online")
}
func (server *Server) SetUpRouter() {
	server.Router = server.IrisApp.Party("/api/v1")
}

func (server *Server) StartListening() {
	server.IrisApp.Run(iris.Addr(":"+server.Port), iris.WithoutBodyConsumptionOnUnmarshal)
}
