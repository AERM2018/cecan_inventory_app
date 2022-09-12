package test

import (
	"cecan_inventory/infrastructure/config"
	"testing"
)

func teardown() {
	server := config.Server{}
	server.New()
	// storage.PruneData(server.DbPsql)
}
func TestServer(t *testing.T) {
	tests := map[string]func(t *testing.T){
		// Login
		"Login should be failed, not found":                          testNotFoundAuth,
		"Login should be ok":                                         testOkAuth,
		"Login should be failed, email's not found":                  testEmailNotFoundAuth,
		"Login should be failed, password's wrong":                   testPasswordWrongAuth,
		"Sign up should be ok":                                       testSignUpOk,
		"Sign up should not be done, wrong role":                     testSignUpRoleWrong,
		"Sign up should not be done, user not valid":                 testSignUpUserWrong,
		"Sign up should not be done, email already in use":           testSignUpEmailUsed,
		"Sign up should not be done, password is short":              testSignUpSmallPassword,
		"Sign up should not be done, email is not valid":             testSignUpEmailNotValid,
		"Token should be refreshed":                                  testRefreshTokenOk,
		"Token should not be refreshed, token is missing":            testRefreshTokenNoAuthHeader,
		"Token should not be refreshed, token is invalid":            testRefreshTokenWithInvalidToken,
		"Medicine should be created":                                 testCreateMedicineOk,
		"Medicine should not be created, it already exists":          testCreateMedicineRepeated,
		"Medicine should not be created, name repeated":              testCreateMedicineNameRepeated,
		"Medicine catalog should have more than 0 elements":          testGetMedicineCatalogOk,
		"Medicine should not be updated, key repeated":               testUpdateMedicineKeyRepeated,
		"Medicine should not be updated, name repeated":              testUpdateMedicineNameRepeated,
		"Medicine should be deleted":                                 testDeleteMedicineOk,
		"Medicine should not be deleted, medicine's not found":       testDeleteMedicineNotFound,
		"Medicine should not be reactivate, medicine is not deleted": testReactivateMedicineNoDeleted,
		"Pharmacy stock should be created":                           testCreatePhStockOk,
		"Pharmacy stock should not be created, medicine no found":    testCreatePhStockMedicineNoFound,
	}
	for name, tt := range tests {
		t.Run(name, tt)
	}
	teardown()
}
