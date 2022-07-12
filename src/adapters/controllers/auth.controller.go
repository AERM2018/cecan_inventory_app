package controllers

import (
	"cecan_inventory/src/adapters/helpers"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type AuthController struct {
	UserDataSource datasources.UserDataSource
}
type AccessCredentials struct{
	Emial string `json:"email"`
	Password string `json:"password"`
}
func (controller AuthController) Login(ctx iris.Context) {
	credentials := AccessCredentials{}
	ctx.ReadBody(&credentials)
	fmt.Println(credentials.Emial)
	user, err := controller.UserDataSource.GetUserByEmail(credentials.Emial)
	if(err != nil){
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.PrepareAndSendMessageResponse(ctx,iris.StatusNotFound,nil,"Invalid credentials.")
			return
		}
	}
	helpers.PrepareAndSendDataResponse(ctx,iris.StatusOK,iris.Map{"user":user})
}
