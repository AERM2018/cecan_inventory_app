package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"cecan_inventory/infrastructure/config"
	"fmt"
	"testing"
	"time"

	"github.com/kataras/iris/v12/httptest"
)

var tokenClaims = models.AuthClaims{Id: "CAN102212", Role: "Admin", FullName: "CECAN ADMIN"}
var token = mocks.GetTokenMock(tokenClaims)

// START create pharmacy stock test cases
func testCreatePhStockOk(t *testing.T) {
	pharmacyStock := mocks.GetPharmacyStockMock(models.Medicine{})
	server := config.Server{}
	app := server.New()
	e := httptest.New(t, app)
	res := e.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", pharmacyStock.MedicineKey)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusCreated)

	res.JSON().Object().Value("data").Object().Value("stock").Object().NotEmpty()
}

func testCreatePhStockMedicineNoFound(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMock()
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	server := config.Server{}
	app := server.New()
	e := httptest.New(t, app)
	res := e.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("El medicamento con clave: %v no existe.", medicineMockSeed.Key))
}

func testCreatePhStockMedicineInactive(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMockSeed()[0]
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	server := config.Server{}
	app := server.New()
	e := httptest.New(t, app)
	// Delete medicine to make it inactive
	e.DELETE(fmt.Sprintf("/api/v1/medicines/%v", medicineMockSeed.Key)).WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).Expect().Status(httptest.StatusNoContent)
	// Insert stock with an inactive medicine
	res := e.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("No se pudó ingresar el stock a farmacia del medicamento con clave: %v debido a que esta inactivo, activelo y vuelvalo a intentar.", medicineMockSeed.Key))
	// reactivate medicine
	e.PUT(fmt.Sprintf("/api/v1/medicines/%v/reactivate", medicineMockSeed.Key)).WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).Expect().Status(httptest.StatusOK)
}

func testCreatePhStockWrongRole(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMockSeed()[0]
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	server := config.Server{}
	claimsWrongRole := tokenClaims
	claimsWrongRole.Role = "Medico"
	newToken := mocks.GetTokenMock(claimsWrongRole)
	app := server.New()
	e := httptest.New(t, app)
	res := e.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", newToken)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusForbidden)

	res.JSON().Object().Value("error").Equal("Acción denegada, no cuenta con los permisos necesarios.")
}

func testCreatePhStockWrongFields(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMockSeed()[0]
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	pharmacyStock.Pieces = 0
	pharmacyStock.ExpiresAt = time.Now().Add(-time.Hour * 48)
	pharmacyStock.LotNumber = ""
	server := config.Server{}
	app := server.New()
	e := httptest.New(t, app)
	res := e.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusBadRequest)

	jsonObj := res.JSON().Object()
	jsonObj.Value("error").Object().Value("errors").Object().Keys().Length().Equal(3)
	jsonObj.Value("error").Object().Value("errors").Object().Value("ExpiresAt").String().Contains("failed on the 'gttoday' tag")
	jsonObj.Value("error").Object().Value("errors").Object().Value("Pieces").String().Contains("Error:Field validation for 'Pieces' failed")
	jsonObj.Value("error").Object().Value("errors").Object().Value("LotNumber").String().Contains("Error:Field validation for 'LotNumber' failed")
}

// END create pharmacy stock test cases
