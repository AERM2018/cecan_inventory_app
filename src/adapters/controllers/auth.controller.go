package controllers

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	usecases "cecan_inventory/domain/useCases"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"

	"strings"

	"github.com/kataras/iris/v12"
)

type AuthController struct {
	UserDataSource               datasources.UserDataSource
	PasswordResetCodesDataSource datasources.PasswordResetCodesDataSource
	Interator                    usecases.AuthInteractor
}

func (controller *AuthController) New(userDatasource datasources.UserDataSource) {
	controller.UserDataSource = userDatasource
	controller.Interator = usecases.AuthInteractor{UserDataSource: userDatasource, PasswordResetCodeDataSource: controller.PasswordResetCodesDataSource}
}

func (controller AuthController) Login(ctx iris.Context) {
	credentials := models.AccessCredentials{}
	ctx.ReadBody(&credentials)
	res := controller.Interator.LoginUser(credentials)
	if res.StatusCode != iris.StatusOK {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller AuthController) SignUp(ctx iris.Context) {
	var newUser models.User
	ctx.ReadBody(&newUser)
	res := controller.Interator.SignUpUser(newUser)
	if res.Err != nil {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller AuthController) RenewToken(ctx iris.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	res := controller.Interator.RenewToken(token)
	if res.Err != nil {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendDataResponse(ctx, res)
}

func (controller AuthController) GeneratePasswordResetCode(ctx iris.Context) {
	var passwordResetInfo models.AccessCredentialsRestart
	bodyreader.ReadBodyAsJson(ctx, &passwordResetInfo, true)
	res := controller.Interator.GeneratePasswordResetCode(passwordResetInfo.Email)
	if res.StatusCode != iris.StatusOK {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendMessageResponse(ctx, res)
}

func (controller AuthController) ResetPassword(ctx iris.Context) {
	userId := ctx.Params().GetStringDefault("userId", "")
	hash := ctx.URLParamDefault("hash", "")
	withOldPassword, _ := ctx.URLParamBool("withOldPassword")
	var passwordResetInfo models.AccessCredentialsRestart
	bodyreader.ReadBodyAsJson(ctx, &passwordResetInfo, true)
	res := controller.Interator.ResetPassword(userId, hash, withOldPassword, passwordResetInfo)
	if res.StatusCode != iris.StatusOK {
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	helpers.PrepareAndSendMessageResponse(ctx, res)
}
