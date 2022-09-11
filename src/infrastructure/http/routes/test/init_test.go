package test

import (
	"cecan_inventory/infrastructure/config"
	"cecan_inventory/infrastructure/storage"
	"testing"
)

func teardown() {
	server := config.Server{}
	server.New()
	storage.PruneData(server.DbPsql)
}
func TestServer(t *testing.T) {
	tests := map[string]func(t *testing.T){
		"Login should be failed, not found":                          testNotFoundAuth,
		"Login should be ok":                                         testOkAuth,
		"Login should be failed, email's not found":                  testEmailNotFoundAuth,
		"Login should be failed, password's wrong":                   testPasswordWrongAuth,
		"Medicine should be created":                                 testCreateMedicineOk,
		"Medicine should not be created, it already exists":          testCreateMedicineRepeated,
		"Medicine should not be created, name repeated":              testCreateMedicineNameRepeated,
		"Medicine catalog should have more than 0 elements":          testGetMedicineCatalogOk,
		"Medicine should not be updated, key repeated":               testUpdateMedicineKeyRepeated,
		"Medicine should not be updated, name repeated":              testUpdateMedicineNameRepeated,
		"Medicine should be deleted":                                 testDeleteMedicineOk,
		"Medicine should not be deleted, medicine's not found":       testDeleteMedicineNotFound,
		"Medicine should not be reactivate, medicine is not deleted": testReactivateMedicineNoDeleted,
	}
	for name, tt := range tests {
		t.Run(name, tt)
	}
	teardown()
}
