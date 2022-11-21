package usecases

import (
	"cecan_inventory/domain/models"
	authtoken "cecan_inventory/infrastructure/external/authToken"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"errors"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"gorm.io/gorm"
)

type AuthInteractor struct {
	UserDataSource datasources.UserDataSource
}

func (interacor AuthInteractor) LoginUser(credentials models.AccessCredentials) models.Responser {
	var (
		user models.User
		err  error
	)
	user, err = interacor.UserDataSource.GetUserByEmailOrId(credentials.Username)
	// User with that email or ID doesn't exist
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Responser{
				StatusCode: iris.StatusNotFound,
				Message:    "Invalid credentials.",
				Err:        nil,
				Data:       nil,
			}
		}
	}
	// The password typed is wrong
	if isCorrectPassword := user.CheckPassword(credentials.Password); !isCorrectPassword {
		return models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    "Invalid credentials.",
			Err:        nil,
			Data:       nil,
		}
	}
	// Generate jwt token
	claims := models.AuthClaims{
		Id:       user.Id,
		Role:     user.Role.Name,
		FullName: user.Name + " " + user.Surname}
	token, err := authtoken.GenerateJWT(claims)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusUnauthorized,
			Message:    "The token couldn't be generated, try later.",
			Err:        nil,
			Data:       nil,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Message:    "",
		Err:        nil,
		Data:       iris.Map{"user": user.WithoutPassword(), "token": token},
	}
}

func (interactor AuthInteractor) SignUpUser(user models.User) models.Responser {
	newUserRecord, err := interactor.UserDataSource.CreateUser(user)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    err.Error(),
			Err:        err,
			Data:       nil,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Message:    "",
		Err:        nil,
		Data:       iris.Map{"user": newUserRecord.WithoutPassword()},
	}
}

func (interactor AuthInteractor) RenewToken(token string) models.Responser {
	verifiedToken, err := jwt.Verify(jwt.HS256, []byte(os.Getenv("JWTSECRET")), []byte(token))
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusUnauthorized,
			Message:    "Invalid token!",
		}
	}
	claims := models.AuthClaims{}
	verifiedToken.Claims(&claims)
	token, errToken := authtoken.GenerateJWT(claims)
	if errToken != nil {
		return models.Responser{
			StatusCode: iris.StatusUnauthorized,
			Message:    "The token couldn't be generated, try later.",
			Err:        nil,
			Data:       nil,
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Message:    "",
		Err:        nil,
		Data:       iris.Map{"token": token},
	}
}
