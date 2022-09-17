package test

import (
	"cecan_inventory/domain/mocks"
	"fmt"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

var medicine = mocks.GetMedicineMock()
var medineMockCreatedByTestKey string
var preDefinedlength = 3 // two added when seed are run and one aded during test

// START ----- Create medicine test templates
func testCreateMedicineOk(t *testing.T) {
	medicine = mocks.GetMedicineMock()
	medineMockCreatedByTestKey = medicine.Key
	res := HttpTester.
		POST("/api/v1/medicines").
		WithJSON(medicine).
		Expect().Status(httptest.StatusCreated)
	res.JSON().Object().ContainsKey("data")
	res.JSON().Object().Value("data").Object().ContainsKey("medicine")
}

func testCreateMedicineRepeated(t *testing.T) {
	medicineRepeated := mocks.GetMedicineMockSeed()[0]
	res := HttpTester.
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
	medicine = mocks.GetMedicineMock()
	medicine.Name = mocks.GetMedicineMockSeed()[0].Name
	res := HttpTester.
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
	res := HttpTester.
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
	medicine = mocks.GetMedicineMockSeed()[0]
	keyToUpdate := medicine.Key
	medicine.Key = mocks.GetMedicineMockSeed()[1].Key
	res := HttpTester.
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
	medicine = mocks.GetMedicineMockSeed()[0]
	keyToUpdate := medicine.Key
	medicine.Name = mocks.GetMedicineMockSeed()[1].Name
	res := HttpTester.
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
	keyToDelete := mocks.GetMedicineMockSeed()[0].Key
	HttpTester.
		DELETE(fmt.Sprintf("/api/v1/medicines/%s", keyToDelete)).
		Expect().Status(httptest.StatusNoContent)

	res := HttpTester.
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
	keyToDelete := mocks.GetMedicineMock().Key
	res := HttpTester.
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
	keyToReactivate := mocks.GetMedicineMockSeed()[0].Key
	res := HttpTester.
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
