package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"fmt"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func testCreatePrescriptionOk(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	prescriptionMock := mocks.GetPrescriptionMock()
	res := httpTester.POST("/api/v1/prescriptions").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(prescriptionMock).
		Expect().Status(httptest.StatusCreated)

	res.JSON().Object().Value("data").Object().Schema(models.PrescriptionDetialed{})
}

func testCreatePrescriptionWrongRole(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	user := getUserByRoleName("Estandar")
	userClaims := models.AuthClaims{Id: user.Id, Role: user.Role.Name, FullName: user.Name + user.Surname}
	userToken := mocks.GetTokenMock(userClaims)
	prescriptionMock := mocks.GetPrescriptionMock()
	res := httpTester.POST("/api/v1/prescriptions").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", userToken)).
		WithJSON(prescriptionMock).
		Expect().Status(httptest.StatusForbidden)

	res.JSON().Object().Value("error").Equal("Acción denegada, no cuenta con los permisos necesarios.")
}

func testCreatePrescriptionMedicineNotFound(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	prescriptionMock := mocks.GetPrescriptionMock()
	prescriptionMock.Medicines = append(
		prescriptionMock.Medicines,
		models.PrescriptionsMedicines{PrescriptionId: prescriptionMock.Id, MedicineKey: mocks.GetMedicineMock().Key})
	res := httpTester.POST("/api/v1/prescriptions").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(prescriptionMock).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Equal("No se pudo crear la receta debido a que no se pudo asignar los medicamentos a la misma.")
}

func testCreatePrescriptionBadStruct(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	prescriptionMock := models.PrescriptionDetialed{}
	res := httpTester.POST("/api/v1/prescriptions").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		WithJSON(prescriptionMock).
		Expect().Status(httptest.StatusBadRequest)

	jsonObj := res.JSON().Object().Value("error").Object().Value("errors")
	jsonObj.Object().Value("instructions").Equal("cannot be blank")
	jsonObj.Object().Value("patient_name").Equal("cannot be blank.")
	jsonObj.Object().Value("medicines").Equal("required key is missing")
}
