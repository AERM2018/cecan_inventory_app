package middlewares

import (
	"cecan_inventory/src/adapters/helpers"
	"cecan_inventory/src/domain/models"
	"fmt"
	"os"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func VerifyJWT(ctx iris.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	fmt.Println(token)
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