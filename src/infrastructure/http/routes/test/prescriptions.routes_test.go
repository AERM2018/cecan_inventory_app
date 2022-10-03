package test

import (
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"fmt"
	"testing"

	"github.com/icrowley/fake"
	"github.com/kataras/iris/v12/httptest"
)

// START create prescription test cases
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

// END create prescription test cases

// START get prescriptions test cases
func testGetPrescriptionsOk(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	res := httpTester.GET("/api/v1/prescriptions").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusOK)
	jsonObj := res.JSON().Object().Value("data").Object().Value("prescriptions")
	jsonObj.Array().Length().Ge(1)
	jsonObj.Array().First().Object().NotContainsKey("medicine")
}

func testGetPrescriptionById(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	prescription := mocks.GetPrescriptionMockSeed()[0]
	res := httpTester.GET("/api/v1/prescriptions/{id}").
		WithPath("id", prescription.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusOK)
	jsonObj := res.JSON().Object().Value("data").Object().Value("prescription")
	jsonObj.Object().Schema(&models.PrescriptionDetialed{})
}

func testGetPrescriptionByUserId(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	doctorUser := getUserByRoleName("medico")
	res := httpTester.GET("/api/v1/prescriptions").
		WithQuery("user_id", doctorUser.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusOK)
	jsonObj := res.JSON().Object().Value("data").Object().Value("prescriptions")
	doctorPrescriptios := jsonObj.Array().Iter()
	for _, val := range doctorPrescriptios {
		val.Object().Value("user_id").Equal(doctorUser.Id)
	}
}

// END get prescriptions test cases

// START update prescription test cases
func testUpdateBasicInfoFromPrescription(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	// Change token info, set doctor user info
	doctorUser := getUserByRoleName("medico")
	doctorUserClaims := models.AuthClaims{Id: doctorUser.Id, Role: doctorUser.Role.Name, FullName: doctorUser.Name + doctorUser.Surname}
	doctorToken := mocks.GetTokenMock(doctorUserClaims)
	// Generate fake description and patient name
	prescription := mocks.GetPrescriptionMockSeed()[0]
	fakeInstructions := mocks.GetPrescriptionMock().Instructions
	fakePatientName := mocks.GetPrescriptionMock().PatientName
	prescription.Instructions = fakeInstructions
	prescription.PatientName = fakePatientName
	res := httpTester.PUT("/api/v1/prescriptions/{id}").
		WithPath("id", prescription.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", doctorToken)).
		WithJSON(prescription).
		Expect().Status(httptest.StatusOK)
	jsonObj := res.JSON().Object().Value("data").Object().Value("prescription")
	jsonObj.Object().Value("instructions").Equal(fakeInstructions)
	jsonObj.Object().Value("patient_name").Equal(fakePatientName)
}

func testUpdatePrescriptionNoCreatorUser(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	// Change token info, set doctor user info
	doctorUser := mocks.GetUserMockSeed(mocks.GetRolesMock("medico")[0].Id.String())
	doctorUserClaims := models.AuthClaims{Id: doctorUser.Id, Role: "medico", FullName: doctorUser.Name + doctorUser.Surname}
	doctorToken := mocks.GetTokenMock(doctorUserClaims)
	// Generate fake description and patient name
	prescription := mocks.GetPrescriptionMockSeed()[0]
	fakeInstructions := mocks.GetPrescriptionMock().Instructions
	prescription.Instructions = fakeInstructions
	res := httpTester.PUT("/api/v1/prescriptions/{id}").
		WithPath("id", prescription.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", doctorToken)).
		WithJSON(prescription).
		Expect().Status(httptest.StatusForbidden)
	res.JSON().Object().Value("error").Equal("Solo el creador de la receta está permitido a actualizarla/borrarla.")
}

func testUpdateMedicinesFromPrescription(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	// Change token info, set doctor user info
	doctorUser := getUserByRoleName("medico")
	doctorUserClaims := models.AuthClaims{Id: doctorUser.Id, Role: doctorUser.Role.Name, FullName: doctorUser.Name + doctorUser.Surname}
	doctorToken := mocks.GetTokenMock(doctorUserClaims)
	// Generate fake description and patient name
	prescription := mocks.GetPrescriptionMockSeed()[0]
	medicineQtyUpdated := prescription.Medicines[0].Pieces + 2
	prescription.Medicines[0].Pieces = medicineQtyUpdated
	res := httpTester.PUT("/api/v1/prescriptions/{id}").
		WithPath("id", prescription.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", doctorToken)).
		WithJSON(prescription).
		Expect().Status(httptest.StatusOK)
	jsonObj := res.JSON().Object().Value("data").Object().Value("prescription")
	jsonObj.Object().Value("medicines").Array().First().Object().Value("pieces").Equal(medicineQtyUpdated)
}

func testCompletePrescription(t *testing.T) {
	var medicineSuppliment []models.PrescriptionsMedicinesToComplete
	httpTester := httptest.New(t, IrisApp)
	prescription := mocks.GetPrescriptionMockSeed()[0]
	// Change token info, set doctor user info
	pharmacyUser := getUserByRoleName("farmacia")
	PharmacyUserClaims := models.AuthClaims{Id: pharmacyUser.Id, Role: pharmacyUser.Role.Name, FullName: pharmacyUser.Name + pharmacyUser.Surname}
	pharmacyToken := mocks.GetTokenMock(PharmacyUserClaims)
	// Supplie medicine from prescription
	for _, medicine := range prescription.Medicines {
		medicineSuppliment = append(
			medicineSuppliment,
			models.PrescriptionsMedicinesToComplete{
				MedicineKey:    medicine.MedicineKey,
				PiecesSupplied: medicine.Pieces,
			})
	}
	supplimentOfPrescription := models.PrescriptionToComplete{
		Observations: fake.Paragraph(),
		Medicines:    medicineSuppliment,
	}
	res := httpTester.PUT("/api/v1/prescriptions/{id}/complete").
		WithPath("id", prescription.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", pharmacyToken)).
		WithJSON(supplimentOfPrescription).
		Expect().Status(httptest.StatusOK)
	jsonObj := res.JSON().Object().Value("data").Object().Value("prescription")
	jsonObj.Object().Schema(models.PrescriptionDetialed{})
	jsonObj.Object().Value("prescription_status").Object().Value("name").Equal("Completada")
	prescriptionMedicines := jsonObj.Object().Value("medicines").Array().Iter()

	for i, medicine := range prescriptionMedicines {
		medicine.Object().Value("pieces_supplied").Number().Equal(medicineSuppliment[i].PiecesSupplied)
	}
}

// TODO: Update a prescription adding a medicine to the prescription medicines list
func testUpdatePrescriptionAddMedicine(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	// Change token info, set doctor user info
	doctorUser := getUserByRoleName("medico")
	doctorUserClaims := models.AuthClaims{Id: doctorUser.Id, Role: doctorUser.Role.Name, FullName: doctorUser.Name + doctorUser.Surname}
	doctorToken := mocks.GetTokenMock(doctorUserClaims)
	// Generate fake description and patient name
	prescription := mocks.GetPrescriptionMockSeed()[0]
	medicineListUpdated := append(prescription.Medicines, models.PrescriptionsMedicines{
		MedicineKey: mocks.GetMedicineMockSeed()[1].Key,
		Pieces:      1,
	})
	prescription.Medicines = medicineListUpdated
	res := httpTester.PUT("/api/v1/prescriptions/{id}").
		WithPath("id", prescription.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", doctorToken)).
		WithJSON(prescription).
		Expect().Status(httptest.StatusOK)
	jsonObj := res.JSON().Object().Value("data").Object().Value("prescription")
	jsonObj.Object().Value("medicines").Array().Length().Equal(len(medicineListUpdated))
}

// END update prescription test cases

// START delete prescription test cases

func testDeletePrescriptionOk(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	prescription := mocks.GetPrescriptionMockSeed()[1]
	// Change token info, set doctor user info
	doctorUser := getUserByRoleName("medico")
	doctorUserClaims := models.AuthClaims{Id: doctorUser.Id, Role: doctorUser.Role.Name, FullName: doctorUser.Name + doctorUser.Surname}
	doctorToken := mocks.GetTokenMock(doctorUserClaims)
	httpTester.DELETE("/api/v1/prescriptions/{id}").
		WithPath("id", prescription.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", doctorToken)).
		Expect().Status(httptest.StatusNoContent)
}

func testDeletePrescriptionNoPendingStatus(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	prescription := mocks.GetPrescriptionMockSeed()[2]
	// Change token info, set doctor user info
	doctorUser := getUserByRoleName("medico")
	doctorUserClaims := models.AuthClaims{Id: doctorUser.Id, Role: doctorUser.Role.Name, FullName: doctorUser.Name + doctorUser.Surname}
	doctorToken := mocks.GetTokenMock(doctorUserClaims)
	res := httpTester.DELETE("/api/v1/prescriptions/{id}").
		WithPath("id", prescription.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", doctorToken)).
		Expect().Status(httptest.StatusBadRequest)

	res.JSON().Object().Value("error").Equal("No se pudó completar la acción, la receta no tiene un estado: pendiente")
}

func testDeletePrescriptionNoSameCreator(t *testing.T) {
	httpTester := httptest.New(t, IrisApp)
	prescription := mocks.GetPrescriptionMockSeed()[0]
	res := httpTester.DELETE("/api/v1/prescriptions/{id}").
		WithPath("id", prescription.Id).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().Status(httptest.StatusForbidden)

	res.JSON().Object().Value("error").Equal("Solo el creador de la receta está permitido a actualizarla/borrarla.")
}

// END delete prescription test cases
