package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/infrastructure/config"
	"fmt"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

var pharmacyStock = mocks.GetPharmacyStockMock()

func testCreatePhStockOk(t *testing.T) {
	server := config.Server{}
	app := server.New()
	e := httptest.New(t, app)
	res := e.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", pharmacyStock.MedicineKey)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusCreated)

	res.JSON().Object().Value("data").Object().Value("stock")
}

func testCreatePhStockMedicineNoFound(t *testing.T) {
	server := config.Server{}
	app := server.New()
	invalidPharmacyStock := pharmacyStock
	invalidPharmacyStock.MedicineKey = mocks.GetMedicineMock().Key
	e := httptest.New(t, app)
	res := e.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", invalidPharmacyStock.MedicineKey)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("El medicamento con clave: %v no existe.", invalidPharmacyStock.MedicineKey))
}
