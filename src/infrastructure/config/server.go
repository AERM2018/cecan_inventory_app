package config

import (
	"fmt"
	"path"
	"runtime"

	"cecan_inventory/src/infrastructure/storage"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type Server struct {
	IrisApp *iris.Application
	DbPsql  *gorm.DB
}

func (server *Server) New() {
	server.IrisApp = iris.New()
	_, filename, _, _ := runtime.Caller(0)
	envPath := path.Join(path.Dir(filename), "../../../.env")
	err := godotenv.Load(envPath)
	if err != nil {
		panic(err)
	}
}

func (server *Server) ConnectDatabase() {
	var errPsql error
	server.DbPsql, errPsql = storage.Connect()
	if errPsql != nil {
		fmt.Println(errPsql)
	}
	fmt.Println("PSQL online")
}
func setUpRoutes(server *Server) {

}
