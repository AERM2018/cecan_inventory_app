package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	authtoken "cecan_inventory/infrastructure/external/authToken"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"github.com/kataras/iris/v12/httptest"
)

var user = mocks.GetUserMock("", 10)

func testOkAuth(t *testing.T) {
	userMockSeed := mocks.GetUserMockSeed("")
	credentials := models.AccessCredentials{Email: userMockSeed.Email, Password: userMockSeed.Password}
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/login").
		WithJSON(credentials).
		Expect().Status(httptest.StatusOK)

	res.JSON().Object().ContainsKey("data")
	res.JSON().Object().Value("data").Object().ContainsKey("user")
	res.JSON().Object().Value("data").Object().ContainsKey("token").Value("token").NotEqual("")
}

func testNotFoundAuth(t *testing.T) {
	credentials := models.AccessCredentials{Email: fake.EmailAddress(), Password: fake.Password(8, 10, true, true, false)}
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/login").
		WithJSON(credentials).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().ValueEqual("error", "Invalid credentials.")
}

func testEmailNotFoundAuth(t *testing.T) {
	credentials := models.AccessCredentials{Email: fake.EmailAddress(), Password: user.Password}
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/login").
		WithJSON(credentials).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().ValueEqual("error", "Invalid credentials.")
}

func testPasswordWrongAuth(t *testing.T) {
	credentials := models.AccessCredentials{Email: user.Email, Password: fake.Password(8, 10, true, true, false)}
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/login").
		WithJSON(credentials).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().ValueEqual("error", "Invalid credentials.")
}

// START test sign up

func testSignUpOk(t *testing.T) {
	user = mocks.GetUserMock(mocks.GetRolesMock()[0].Id.String(), 10)
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusOK)

	res.JSON().Object().Value("data").Object().Value("user").Object().Schema(user)
}

func testSignUpRoleWrong(t *testing.T) {
	user = mocks.GetUserMock(uuid.NewString(), 10)
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("El rol con id: %v no existe.", user.RoleId))
}

func testSignUpUserWrong(t *testing.T) {
	user = mocks.GetUserMock("", 10)
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error")
}

func testSignUpEmailUsed(t *testing.T) {
	user = mocks.GetUserMockSeed(mocks.GetRolesMock()[0].Id.String())
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("El email %v ya está siendo usado por otro usuario.", user.Email))
}

func testSignUpSmallPassword(t *testing.T) {
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMock(role.Id.String(), 5)
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Object().Value("errors").Object().Value("password").String().Contains("the length must be between")
}

func testSignUpEmailNotValid(t *testing.T) {
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMock(role.Id.String(), 9)
	user.Email = "xxxxxx.com"
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/signup").
		WithJSON(user).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Object().Value("errors").Object().Value("email").Equal("must be a valid email address.")
}

// END test sign up

// START test refresh token
func testRefreshTokenOk(t *testing.T) {
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMockSeed(role.Id.String())
	claims := models.AuthClaims{Id: user.Id, Role: role.Id.String(), FullName: user.Name + user.Surname}
	token, _ := authtoken.GenerateJWT(claims)
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/renew").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusOK)

	res.JSON().Object().Value("data").Object().Value("token")
}

func testRefreshTokenNoAuthHeader(t *testing.T) {
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMockSeed(role.Id.String())
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/renew").
		Expect().Status(httptest.StatusUnauthorized)

	res.JSON().Object().Value("error").Equal("Invalid token!")
}

func testRefreshTokenWithInvalidToken(t *testing.T) {
	var role models.Role
	server.DbPsql.First(&role)
	user = mocks.GetUserMockSeed(role.Id.String())
	claims := models.AuthClaims{Id: user.Id, Role: role.Id.String(), FullName: user.Name + user.Surname}
	token, _ := authtoken.GenerateJWT(claims)
	token += "a"
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.
		POST("/api/v1/auth/renew").
		Expect().Status(httptest.StatusUnauthorized)

	res.JSON().Object().Value("error").Equal("Invalid token!")
}

// END test refresh token
