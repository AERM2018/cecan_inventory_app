package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/infrastructure/config"
	"fmt"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

var medicine = mocks.GetMedicineMock()
var medineMockCreatedByTestKey string
var preDefinedlength = 3 // two added when seed are run and one aded during test

// START ----- Create medicine test templates
func testCreateMedicineOk(t *testing.T) {
	server := config.Server{}
	app := server.New()
	medicine = mocks.GetMedicineMock()
	medineMockCreatedByTestKey = medicine.Key
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/medicines").
		WithJSON(medicine).
		Expect().Status(httptest.StatusCreated)
	res.JSON().Object().ContainsKey("data")
	res.JSON().Object().Value("data").Object().ContainsKey("medicine")
}

func testCreateMedicineRepeated(t *testing.T) {
	server := config.Server{}
	app := server.New()
	medicineRepeated := mocks.GetMedicineMockSeed()[0]
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/medicines").
		WithJSON(medicineRepeated).
		Expect().Status(httptest.StatusBadRequest)
	res.
		JSON().
		Object().
		ContainsKey("error").
		Value("error").
		Equal(fmt.Sprintf("El medicamento con clave: %s ya existe.", medicineRepeated.Key))
}

func testCreateMedicineNameRepeated(t *testing.T) {
	server := config.Server{}
	app := server.New()
	medicine = mocks.GetMedicineMock()
	medicine.Name = mocks.GetMedicineMockSeed()[0].Name
	e := httptest.New(t, app)
	res := e.
		POST("/api/v1/medicines").
		WithJSON(medicine).
		Expect().Status(httptest.StatusBadRequest)
	res.
		JSON().
		Object().
		ContainsKey("error").
		Value("error").
		String().Contains(fmt.Sprintf("Ya existe un medicamento con el nombre: %s.", medicine.Name))
}

// END ----- Create medicine test templates

// START ----- Get medicines test templates
func testGetMedicineCatalogOk(t *testing.T) {
	server := config.Server{}
	app := server.New()
	e := httptest.New(t, app)
	res := e.
		GET("/api/v1/medicines").
		Expect().Status(httptest.StatusOK)
	res.
		JSON().
		Object().
		ContainsKey("data").
		Value("data").
		Object().
		ContainsKey("medicines").
		Value("medicines").
		Array().
		Length().
		Gt(0)
}

// END ----- Get medicines test templates

// START ----- Update medicine test templates
func testUpdateMedicineKeyRepeated(t *testing.T) {
	server := config.Server{}
	app := server.New()
	medicine = mocks.GetMedicineMockSeed()[0]
	keyToUpdate := medicine.Key
	medicine.Key = mocks.GetMedicineMockSeed()[1].Key
	e := httptest.New(t, app)
	res := e.
		PUT(fmt.Sprintf("/api/v1/medicines/%s", keyToUpdate)).
		WithJSON(medicine).
		Expect().Status(httptest.StatusBadRequest)
	res.
		JSON().
		Object().
		ContainsKey("error").
		ValueEqual("error", fmt.Sprintf("No se actualizó el medicamento debido a que ya existe un medicamento con la clave: %v.", medicine.Key))
}

func testUpdateMedicineNameRepeated(t *testing.T) {
	server := config.Server{}
	app := server.New()
	medicine = mocks.GetMedicineMockSeed()[0]
	keyToUpdate := medicine.Key
	medicine.Name = mocks.GetMedicineMockSeed()[1].Name
	e := httptest.New(t, app)
	res := e.
		PUT(fmt.Sprintf("/api/v1/medicines/%s", keyToUpdate)).
		WithJSON(medicine).
		Expect().Status(httptest.StatusBadRequest)
	res.
		JSON().
		Object().
		ContainsKey("error").
		Value("error").
		String().
		Contains(fmt.Sprintf("No se actualizó el medicamento debido a que ya existe un medicamento con el nombre: %v.", medicine.Name))
}

// END ----- Update medicines test templates

// START ----- Delete medicine test templates

func testDeleteMedicineOk(t *testing.T) {
	server := config.Server{}
	app := server.New()
	keyToDelete := mocks.GetMedicineMockSeed()[0].Key
	e := httptest.New(t, app)
	e.
		DELETE(fmt.Sprintf("/api/v1/medicines/%s", keyToDelete)).
		Expect().Status(httptest.StatusNoContent)

	res := e.
		PUT(fmt.Sprintf("/api/v1/medicines/%s/reactivate", keyToDelete)).
		Expect().Status(httptest.StatusOK)
	res.
		JSON().
		Object().
		ContainsKey("data").
		Value("data").
		Object().
		ContainsKey("medicine")
}

func testDeleteMedicineNotFound(t *testing.T) {
	server := config.Server{}
	app := server.New()
	keyToDelete := mocks.GetMedicineMock().Key
	e := httptest.New(t, app)
	res := e.
		DELETE(fmt.Sprintf("/api/v1/medicines/%s", keyToDelete)).
		Expect().Status(httptest.StatusNotFound)
	res.
		JSON().
		Object().
		ContainsKey("error").
		Value("error").
		Equal(fmt.Sprintf("El medicamento con clave: %v no existe.", keyToDelete))
}

func testReactivateMedicineNoDeleted(t *testing.T) {
	server := config.Server{}
	app := server.New()
	keyToReactivate := mocks.GetMedicineMockSeed()[0].Key
	e := httptest.New(t, app)
	res := e.
		PUT(fmt.Sprintf("/api/v1/medicines/%s/reactivate", keyToReactivate)).
		Expect().Status(httptest.StatusBadRequest)
	res.
		JSON().
		Object().
		ContainsKey("error").
		Value("error").
		Equal(fmt.Sprintf("El medicamento con clave: %v no se reactivó debido a que no ha sido eliminado antes.", keyToReactivate))
}

// END ----- Delete medicine test templates
