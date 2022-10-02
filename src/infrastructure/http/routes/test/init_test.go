package test

import (
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/mocks"
	"cecan_inventory/domain/models"
	"cecan_inventory/infrastructure/config"
	"cecan_inventory/infrastructure/storage"
	"strings"
	"testing"

	"github.com/kataras/iris/v12"
)

var (
	server      config.Server
	IrisApp     *iris.Application
	token       string
	tokenClaims models.AuthClaims
	userMocks   []models.User
)

func initServerTester(t *testing.T) {
	server = config.Server{}
	IrisApp = server.New()
	server.DbPsql.Joins("Role").Find(&userMocks)
	userFounds := common.FilterSlice(userMocks, func(i interface{}) bool {
		parsed := i.(models.User)
		return strings.ToLower(parsed.Role.Name) == strings.ToLower("Admin")
	})
	adminUser := userFounds.([]models.User)[0]
	tokenClaims = models.AuthClaims{Id: adminUser.Id, Role: adminUser.Role.Name, FullName: adminUser.Name + adminUser.Surname}
	token = mocks.GetTokenMock(tokenClaims)
}

func getUserByRoleName(name string) models.User {
	usersFound := common.FilterSlice(userMocks, func(i interface{}) bool {
		parsed := i.(models.User)
		return strings.ToLower(parsed.Role.Name) == strings.ToLower(name)
	})
	return usersFound.([]models.User)[0]
}

func teardown() {
	storage.PruneData(server.DbPsql)
}
func TestServer(t *testing.T) {
	initServerTester(t)

	tests := map[string]func(t *testing.T){
		"Login should be failed, not found":                                  testNotFoundAuth,
		"Login should be ok":                                                 testOkAuth,
		"Login should be failed, email's not found":                          testEmailNotFoundAuth,
		"Login should be failed, password's wrong":                           testPasswordWrongAuth,
		"Sign up should be ok":                                               testSignUpOk,
		"Sign up should not be done, wrong role":                             testSignUpRoleWrong,
		"Sign up should not be done, user not valid":                         testSignUpUserWrong,
		"Sign up should not be done, email already in use":                   testSignUpEmailUsed,
		"Sign up should not be done, password is short":                      testSignUpSmallPassword,
		"Sign up should not be done, email is not valid":                     testSignUpEmailNotValid,
		"Token should be refreshed":                                          testRefreshTokenOk,
		"Token should not be refreshed, token is missing":                    testRefreshTokenNoAuthHeader,
		"Token should not be refreshed, token is invalid":                    testRefreshTokenWithInvalidToken,
		"Medicine should be created":                                         testCreateMedicineOk,
		"Medicine should not be created, it already exists":                  testCreateMedicineRepeated,
		"Medicine should not be created, name repeated":                      testCreateMedicineNameRepeated,
		"Medicine catalog should have more than 0 elements":                  testGetMedicineCatalogOk,
		"Medicine should not be updated, key repeated":                       testUpdateMedicineKeyRepeated,
		"Medicine should not be updated, name repeated":                      testUpdateMedicineNameRepeated,
		"Medicine should be deleted":                                         testDeleteMedicineOk,
		"Medicine should not be deleted, medicine's not found":               testDeleteMedicineNotFound,
		"Medicine should not be reactivate, medicine is not deleted":         testReactivateMedicineNoDeleted,
		"Pharmacy stocks should be retrieved":                                testGetPharmacyStocksOk,
		"Pharmacy stocks of a medicine should be retrieved":                  testGetPharmacyStockOfMedicine,
		"Pharmacy stock should be created":                                   testCreatePhStockOk,
		"Pharmacy stock should not be created, medicine no found":            testCreatePhStockMedicineNoFound,
		"Pharmacy stock should not be created, medicine is inactive":         testCreatePhStockMedicineInactive,
		"Pharmacy stock should not be created, user doesn't have rigth role": testCreatePhStockWrongRole,
		"Pharmacy stock should not be created, fail struct validation":       testCreatePhStockWrongFields,
		"Pharmacy stock should be updated":                                   testUpdatePhStockOk,
		// "Pharmacy stock should not be updated, medicine not found":           testUpdatePhStockNotFound,
		"Pharmacy stock should deleted":                              testDeletePhStockOk,
		"Pharmacy stock should not be deleted, wrong role":           testDeletePhStockWrongRole,
		"Pharmacy stock should not be deleted, stock doesnp't exist": testDeletePhStockNotFound,
		"Pharmacy stock should not be deleted, it's already in use":  testDeletePhStockUsed,
		"Prescription should be created":                             testCreatePrescriptionOk,
		"Prescription should not be created, wrong role":             testCreatePrescriptionWrongRole,
		"Prescription should not be created, medicine not found":     testCreatePrescriptionMedicineNotFound,
		"Prescription should not be created, info's missing":         testCreatePrescriptionBadStruct,
		"Prescriptions should retrieved":                             testGetPrescriptionsOk,
		"Prescription by id should  retrieved":                       testGetPrescriptionById,
		"Prescriptions of doctor should  retrieved":                  testGetPrescriptionByUserId,
	}
	for name, tt := range tests {
		t.Run(name, tt)
	}

	teardown()
}
