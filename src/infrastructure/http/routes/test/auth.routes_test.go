package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"cecan_inventory/infrastructure/config"
	authtoken "cecan_inventory/infrastructure/external/authToken"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"github.com/kataras/iris/v12/httptest"
)

var user = mocks.GetUserMock("", 10)

func testOkAuth(t *testing.T) {
	server := config.Server{}
	app := server.New()
	userMockSeed := mocks.GetUserMockSeed("")
	credentials := models.AccessCredentials{Email: userMockSeed.Email, Password: userMockSeed.Password}
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/login").
		WithJSON(credentials).
		Expect().Status(httptest.StatusOK)

	res.JSON().Object().ContainsKey("data")
	res.JSON().Object().Value("data").Object().ContainsKey("user")
	res.JSON().Object().Value("data").Object().ContainsKey("token").Value("token").NotEqual("")
}

func testNotFoundAuth(t *testing.T) {
	server := config.Server{}
	app := server.New()
	credentials := models.AccessCredentials{Email: fake.EmailAddress(), Password: fake.Password(8, 10, true, true, false)}
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/login").
		WithJSON(credentials).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().ValueEqual("error", "Invalid credentials.")
}

func testEmailNotFoundAuth(t *testing.T) {
	server := config.Server{}
	app := server.New()
	credentials := models.AccessCredentials{Email: fake.EmailAddress(), Password: user.Password}
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/login").
		WithJSON(credentials).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().ValueEqual("error", "Invalid credentials.")
}

func testPasswordWrongAuth(t *testing.T) {
	server := config.Server{}
	app := server.New()
	credentials := models.AccessCredentials{Email: user.Email, Password: fake.Password(8, 10, true, true, false)}
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/login").
		WithJSON(credentials).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().ValueEqual("error", "Invalid credentials.")
}

// START test sign up

func testSignUpOk(t *testing.T) {
	server := config.Server{}
	app := server.New()
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMock(role.Id.String(), 10)
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusOK)

	res.JSON().Object().Value("data").Object().Value("user").Object().Schema(user)
}

func testSignUpRoleWrong(t *testing.T) {
	server := config.Server{}
	app := server.New()
	user = mocks.GetUserMock(uuid.NewString(), 10)
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("El rol con id: %v no existe.", user.RoleId))
}

func testSignUpUserWrong(t *testing.T) {
	server := config.Server{}
	app := server.New()
	user = mocks.GetUserMock("", 10)
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error")
}

func testSignUpEmailUsed(t *testing.T) {
	server := config.Server{}
	app := server.New()
	user = mocks.GetUserMockSeed("")
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("El email %v ya está siendo usado por otro usuario.", user.Email))
}

func testSignUpSmallPassword(t *testing.T) {
	server := config.Server{}
	app := server.New()
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMock(role.Id.String(), 5)
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Object().Value("errors").Object().Value("Password").String().Contains("Error:Field validation")
}

func testSignUpEmailNotValid(t *testing.T) {
	server := config.Server{}
	app := server.New()
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMock(role.Id.String(), 9)
	user.Email = "xxxxxx.com"
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Object().Value("errors").Object().Value("Email").String().Contains("Error:Field validation")
}

// END test sign up

// START test refresh token
func testRefreshTokenOk(t *testing.T) {
	server := config.Server{}
	app := server.New()
	e := httptest.New(t, app)
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMockSeed(role.Id.String())
	claims := models.AuthClaims{Id: user.Id, Role: role.Id.String(), FullName: user.Name + user.Surname}
	token, _ := authtoken.GenerateJWT(claims)
	res := e.
		POST("/api/v1/auth/renew").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusOK)

	res.JSON().Object().Value("data").Object().Value("token")
}

func testRefreshTokenNoAuthHeader(t *testing.T) {
	server := config.Server{}
	app := server.New()
	e := httptest.New(t, app)
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMockSeed(role.Id.String())
	res := e.
		POST("/api/v1/auth/renew").
		Expect().Status(httptest.StatusUnauthorized)

	res.JSON().Object().Value("error").Equal("Invalid token!")
}

func testRefreshTokenWithInvalidToken(t *testing.T) {
	server := config.Server{}
	app := server.New()
	e := httptest.New(t, app)
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMockSeed(role.Id.String())
	claims := models.AuthClaims{Id: user.Id, Role: role.Id.String(), FullName: user.Name + user.Surname}
	token, _ := authtoken.GenerateJWT(claims)
	token += "a"
	res := e.
		POST("/api/v1/auth/renew").
		Expect().Status(httptest.StatusUnauthorized)

	res.JSON().Object().Value("error").Equal("Invalid token!")
}

// END test refresh token
