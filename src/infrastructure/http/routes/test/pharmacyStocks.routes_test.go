package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"cecan_inventory/infrastructure/storage/seeds"
	"fmt"
	"testing"
	"time"

	"github.com/kataras/iris/v12/httptest"
)

// START get pharmacy inventory test cases
func testGetPharmacyStocksOk(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.GET("/api/v1/pharmacy_inventory").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusOK)
	jsonObj := res.JSON().Object().Value("data")
	pharmacyStockObj := jsonObj.Object().Value("inventory").Array().Element(0)
	pharmacyStockObj.Object().Value("medicine").Schema(models.Medicine{})
	pharmacyStockObj.Object().Value("pieces_by_semaforization_color")
	pharmacyStockObj.Object().Value("stocks").Array().Length().Ge(0)
	pharmacyStockObj.Object().Value("total_pieces").NotNull()
}

func testGetPharmacyStockOfMedicine(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	medicine := mocks.GetMedicineMockSeed()[0]
	res := httpTester.GET("/api/v1/pharmacy_inventory").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithQuery("medicine_key", medicine.Key).
		Expect().Status(httptest.StatusOK)

	jsonObj := res.JSON().Object().Value("data")
	jsonObj.Object().Value("inventory").Array().Length().Equal(1)
	pharmacyStockObj := jsonObj.Object().Value("inventory").Array().Element(0)
	pharmacyStockObj.Object().Value("medicine").Schema(models.Medicine{})
	pharmacyStockObj.Object().Value("pieces_by_semaforization_color")
	pharmacyStockObj.Object().Value("stocks").Array().Length().Ge(0)
	pharmacyStockObj.Object().Value("total_pieces").NotNull()
}

// END get pharmacy inventory test cases

// START create pharmacy stock test cases
func testCreatePhStockOk(t *testing.T) {
	pharmacyStock := mocks.GetPharmacyStockMock(models.Medicine{})
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", pharmacyStock.MedicineKey)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusCreated)

	res.JSON().Object().Value("data").Object().Value("stock").Object().NotEmpty()
}

func testCreatePhStockMedicineNoFound(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMock()
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("El medicamento con clave: %v no existe.", medicineMockSeed.Key))
}

func testCreatePhStockMedicineInactive(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMockSeed()[0]
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	// Delete medicine to make it inactive
	httpTester := httptest.New(t, IrisApp)
	httpTester.DELETE(fmt.Sprintf("/api/v1/medicines/%v", medicineMockSeed.Key)).WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).Expect().Status(httptest.StatusNoContent)
	// Insert stock with an inactive medicine
	res := httpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(pharmacyStock).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("No se pudó ingresar el stock a farmacia del medicamento con clave: %v debido a que esta inactivo, activelo y vuelvalo a intentar.", medicineMockSeed.Key))
	// reactivate medicine
	httpTester.PUT(fmt.Sprintf("/api/v1/medicines/%v/reactivate", medicineMockSeed.Key)).WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).Expect().Status(httptest.StatusOK)
}

func testCreatePhStockWrongRole(t *testing.T) {
	medicineMockSeed := mocks.GetMedicineMockSeed()[0]
	pharmacyStock := mocks.GetPharmacyStockMock(medicineMockSeed)
	claimsWrongRole := tokenClaims
	claimsWrongRole.Role = "Medico"
	newToken := mocks.GetTokenMock(claimsWrongRole)
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
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
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.POST(fmt.Sprintf("/api/v1/medicines/%v/pharmacy_inventory", medicineMockSeed.Key)).
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
	stock := mocks.GetPharmacyStockMockSeed()[1]
	newStock := mocks.GetPharmacyStockMock(mocks.GetMedicineMockSeed()[0])
	stockToUpdate := models.PharmacyStockToUpdate{
		MedicineKey: newStock.MedicineKey,
		LotNumber:   newStock.LotNumber,
		Pieces:      newStock.Pieces,
		ExpiresAt:   newStock.ExpiresAt,
	}
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.PUT(fmt.Sprintf("/api/v1/pharmacy_inventory/%v", stock.Id)).
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
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.PUT(fmt.Sprintf("/api/v1/pharmacy_inventory/%v", pharmacyStock.Id)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(stockToUpdate).
		Expect().Status(httptest.StatusNotFound)

	res.JSON().Object().Value("error")

}

// END update pharmacy stock test cases

// STRAT delete pharmacy stock test cases
func testDeletePhStockOk(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	pharmacyStock := mocks.GetPharmacyStockMockSeed()[0]
	httpTester.DELETE(fmt.Sprintf("/api/v1/pharmacy_inventory/%v", pharmacyStock.Id)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusNoContent)

	seeds.CreatePharmacyStock(server.DbPsql)
}

func testDeletePhStockWrongRole(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	wrongRoleUser := models.AuthClaims{Id: "CECAN100121", Role: "Medico", FullName: "CECAN TEST"}
	wrongRoleToken := mocks.GetTokenMock(wrongRoleUser)
	pharmacyStock := mocks.GetPharmacyStockMockSeed()[0]
	httpTester.DELETE(fmt.Sprintf("/api/v1/pharmacy_inventory/%v", pharmacyStock.Id)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", wrongRoleToken)).
		Expect().Status(httptest.StatusForbidden)
}

func testDeletePhStockNotFound(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	pharmacyStock := mocks.GetPharmacyStockMock(mocks.GetMedicineMockSeed()[0])
	res := httpTester.DELETE(fmt.Sprintf("/api/v1/pharmacy_inventory/%v", pharmacyStock.Id)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Equal(fmt.Sprintf("El stock de farmacia con id: %v no existe.", pharmacyStock.Id))
}

func testDeletePhStockUsed(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	pharmacyStock := mocks.GetPharmacyStockMockSeed()[1]
	res := httpTester.DELETE(fmt.Sprintf("/api/v1/pharmacy_inventory/%v", pharmacyStock.Id)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Equal("No se puede eliminar un stock de farmacia cuando ha sido utilizado.")
}

// END delete pharmacy stock test cases
