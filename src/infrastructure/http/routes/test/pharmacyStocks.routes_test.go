package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
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
	res := HttpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", pharmacyStock.MedicineKey)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusCreated)

	res.JSON().Object().Value("data").Object().Value("stock").Object().NotEmpty()
}

func testCreatePhStockMedicineNoFound(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMock()
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	res := HttpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("El medicamento con clave: %v no existe.", medicineMockSeed.Key))
}

func testCreatePhStockMedicineInactive(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMockSeed()[0]
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	// Delete medicine to make it inactive
	HttpTester.DELETE(fmt.Sprintf("/api/v1/medicines/%v", medicineMockSeed.Key)).WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).Expect().Status(httptest.StatusNoContent)
	// Insert stock with an inactive medicine
	res := HttpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("No se pudó ingresar el stock a farmacia del medicamento con clave: %v debido a que esta inactivo, activelo y vuelvalo a intentar.", medicineMockSeed.Key))
	// reactivate medicine
	HttpTester.PUT(fmt.Sprintf("/api/v1/medicines/%v/reactivate", medicineMockSeed.Key)).WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).Expect().Status(httptest.StatusOK)
}

func testCreatePhStockWrongRole(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMockSeed()[0]
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	claimsWrongRole := tokenClaims
	claimsWrongRole.Role = "Medico"
	newToken := mocks.GetTokenMock(claimsWrongRole)
	res := HttpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
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
	res := HttpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusBadRequest)

	jsonObj := res.JSON().Object()
	jsonObj.Value("error").Object().Value("errors").Object().Keys().Length().Equal(3)
	jsonObj.Value("error").Object().Value("errors").Object().Value("expires_at").Equal("the date is out of range")
	jsonObj.Value("error").Object().Value("errors").Object().Value("pieces").Equal("must be more than 0 pieces.")
	jsonObj.Value("error").Object().Value("errors").Object().Value("lot_number").Equal("cannot be blank")
}

// END create pharmacy stock test cases

// START update pharmacy stock test cases
func testUpdatePhStockOk(t *testing.T) {
	newStock := mocks.GetPharmacyStockMockSeed()
	stockToUpdate := models.PharmacyStockToUpdate{
		MedicineKey: newStock.MedicineKey,
		LotNumber:   newStock.LotNumber,
		Pieces:      newStock.Pieces,
		ExpiresAt:   newStock.ExpiresAt,
	}
	res := HttpTester.PUT(fmt.Sprintf("/api/v1/pharmacy_inventory/%v", newStock.Id)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(stockToUpdate).
		Expect().Status(httptest.StatusOK)

	phStockInResponse := res.JSON().Object().Value("data").Object().Value("stock").Object()
	phStockInResponse.Value("medicine_key").Equal(stockToUpdate.MedicineKey)
	phStockInResponse.Value("lot_number").Equal(stockToUpdate.LotNumber)
	phStockInResponse.Value("pieces").Equal(stockToUpdate.Pieces)
	phStockInResponse.Value("expires_at").Equal(stockToUpdate.ExpiresAt)
}

func testUpdatePhStockNotFound(t *testing.T) {
	var pharmacyStock models.PharmacyStock
	server.DbPsql.Joins("Medicine").First(&pharmacyStock)
	newStock := mocks.GetPharmacyStockMock(*pharmacyStock.Medicine)
	stockToUpdate := models.PharmacyStockToUpdate{
		MedicineKey: pharmacyStock.MedicineKey,
		LotNumber:   newStock.LotNumber,
		Pieces:      newStock.Pieces,
		ExpiresAt:   newStock.ExpiresAt,
	}
	res := HttpTester.PUT(fmt.Sprintf("/api/v1/pharmacy_inventory/%v", pharmacyStock.Id)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(stockToUpdate).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().Value("error")

}

// END update pharmacy stock test cases
