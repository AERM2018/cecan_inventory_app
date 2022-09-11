package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"cecan_inventory/infrastructure/config"
	"testing"

	"github.com/icrowley/fake"
	"github.com/kataras/iris/v12/httptest"
)

var user = mocks.GetUserMock("")

func testOkAuth(t *testing.T) {
	server := config.Server{}
	app := server.New()
	credentials := models.AccessCredentials{Email: user.Email, Password: user.Password}
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
