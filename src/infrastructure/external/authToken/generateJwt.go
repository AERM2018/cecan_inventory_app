package authtoken

import (
	"cecan_inventory/src/domain/models"
	"os"
	"time"

	"github.com/kataras/iris/v12/middleware/jwt"
)

func GenerateJWT(claims models.AuthClaims) (string, error) {
	signer := jwt.NewSigner(jwt.HS256, os.Getenv("JWTSECRET"), 60*time.Minute)
	token, err := signer.Sign(claims)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
