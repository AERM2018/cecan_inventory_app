package usecases

import (
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	authtoken "cecan_inventory/infrastructure/external/authToken"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/icrowley/fake"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"gorm.io/gorm"
)

type AuthInteractor struct {
	UserDataSource              datasources.UserDataSource
	PasswordResetCodeDataSource datasources.PasswordResetCodesDataSource
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

func (interactor AuthInteractor) GeneratePasswordResetCode(email string) models.Responser {
	randomCode := strings.ToUpper(fake.CharactersN(8))
	userEmailOwner, _ := interactor.UserDataSource.GetUserByEmailOrId(email)
	passwordResetCode := models.PasswordResetCode{
		Code:      randomCode,
		ExpiresAt: time.Now().Add(time.Minute * 3),
	}
	hashVerificationCode := passwordResetCode.HashToken()
	resetPasswordUrl := fmt.Sprintf("%v/auth/password_reset/users/%v?withOldPassword=false&hash=%v",
		os.Getenv("API_URL"),
		userEmailOwner.Id,
		hashVerificationCode,
	)
	errCreatingCode := interactor.PasswordResetCodeDataSource.CreateCode(passwordResetCode)
	if errCreatingCode != nil {
		if errCreatingCode != nil {
			return models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    "Ocurrió un error al crear su código de restablecimiento de contraseña, intentelo mas tarde.",
			}
		}
	}
	err := common.SendPlainEmail(email, "Reestablecimiento de contraseña de CECAN INVENTORY APP\n", "El siguiente código es para el restablecimiento de su contraseña, así que no debe compartirlo con nadie\n Código: "+randomCode)
	if err != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    "Ocurrió un error al enviar el correo de restablecimiento de contraseña, intentelo mas tarde.",
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Message:    "El código fue enviado satisfactoriamente al correo indicado.",
		ExtraInfo: []map[string]interface{}{
			{"url": resetPasswordUrl},
		},
	}
}
func (interactor AuthInteractor) ResetPassword(id string, hash string, withOldPassword bool, passwordReset models.AccessCredentialsRestart) models.Responser {
	if withOldPassword && passwordReset.OldPassword == "" {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    "Acción denegada, enviar la contraseña anterior es requerido.",
		}
	}
	resetCode, errNotFound := interactor.PasswordResetCodeDataSource.GetCode(passwordReset.ResetCode)
	if errors.Is(errNotFound, gorm.ErrRecordNotFound) {
		return models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("El código %v no existe.", passwordReset.ResetCode),
		}
	}
	validToken := resetCode.CheckToken(hash)
	if !validToken {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("El hash para el uso de código %v es incorrecto.", passwordReset.ResetCode),
		}
	}
	if resetCode.IsUsed || resetCode.ExpiresAt.Unix() < time.Now().Unix() {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("El código %v ya fue utilizado o ya expiró.", passwordReset.ResetCode),
		}
	}
	user, _ := interactor.UserDataSource.GetUserByEmailOrId(id)
	user.RestPassword(passwordReset.NewPassword)
	errPasswordUpdate := interactor.UserDataSource.UpdateUserPassword(user)
	interactor.PasswordResetCodeDataSource.UseCode(passwordReset.ResetCode)
	if errPasswordUpdate != nil {
		return models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    "Ocurrió un error al restablecer su contraseña, intentelo mas tarde.",
		}
	}
	return models.Responser{
		StatusCode: iris.StatusOK,
		Message:    "Contraseña restablecida satisfactoriamente.",
	}
}
