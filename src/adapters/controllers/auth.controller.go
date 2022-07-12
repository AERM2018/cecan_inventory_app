package controllers

import (
	datasources "cecan_inventory/src/infrastructure/external/dataSources"
	"fmt"

	"github.com/kataras/iris/v12"
)

type AuthController struct {
	userDataSource datasources.UserDataSource
}

func (controller AuthController) Login(ctx iris.Context) {
	res := controller.userDataSource.Login()
	fmt.Println(res)
	ctx.JSON(iris.Map{
		"ok":   true,
		"data": res,
	})
}
