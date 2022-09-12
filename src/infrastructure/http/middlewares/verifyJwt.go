package middlewares

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	"os"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func VerifyJWT(ctx iris.Context) {
	var token string
	if ctx.GetHeader("Authorization") != "" {
		token = strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	}
	_, err := jwt.Verify(jwt.HS256, []byte(os.Getenv("JWTSECRET")), []byte(token))
	if err != nil {
		res := models.Responser{
			StatusCode: iris.StatusUnauthorized,
			Message:    "Invalid token!",
		}
		helpers.PrepareAndSendMessageResponse(ctx, res)
		return
	}
	ctx.Next()
}
