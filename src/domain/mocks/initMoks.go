package mocks

import (
	"cecan_inventory/domain/models"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/icrowley/fake"
)

func GetRolesMock() []models.Role {
	rolesMocksAsMap := []map[string]string{
		{"id": "7d3f3faa-39e2-4b3c-aaa8-8ca60fa090b4", "name": "Medico"},
		{"id": "6648c302-f487-491d-be57-5f90bbe380c7", "name": "Farmacia"},
		{"id": "5b4e7720-b663-474b-841d-aed481907908", "name": "Almacén"},
		{"id": "6542908e-41ca-4230-8b45-985be3392b5d", "name": "Estandar"},
		{"id": "097305e0-c4a3-4fd6-a929-49cd1aca5e90", "name": "Ingeniero"},
		{"id": "e2660b8f-404c-4baa-b7f5-34e98c602046", "name": "Auditor"},
		{"id": "3c9c3b1c-80c5-43a0-9276-4f88b67a7ac7", "name": "Admin"},
	}
	var rolesMocks []models.Role
	for _, roleMockAsMap := range rolesMocksAsMap {
		idAsUuid, _ := uuid.FromBytes([]byte(roleMockAsMap["id"]))
		fmt.Println(roleMockAsMap["name"])
		rolesMocks = append(rolesMocks, models.Role{Id: &idAsUuid, Name: roleMockAsMap["name"]})
	}
	return rolesMocks
}

func GetUserMockSeed(rolId string) models.User {
	return models.User{
		RoleId:   rolId,
		Password: "Qwerty*123",
		Name:     "CECAN ADMIN",
		Email:    "admin@cecan.com",
	}
}

func GetUserMock(rolId string, minPassLen int) models.User {
	return models.User{
		RoleId:   rolId,
		Password: fake.Password(minPassLen, minPassLen, true, true, false),
		Name:     fake.MaleFirstName(),
		Email:    fake.EmailAddress(),
	}
}

func GetMedicineMock() models.Medicine {
	return models.Medicine{
		Key:  fake.DigitsN(9),
		Name: fake.ProductName(),
	}
}

func GetMedicineMockSeed() []models.Medicine {
	return []models.Medicine{
		{
			Key:  "000000002",
			Name: "Paracetamol",
		}, {
			Key:  "000000009",
			Name: "Naproxeno",
		},
	}
}

func GetPharmacyStockMock() models.PharmacyStock {
	medicine := GetMedicineMockSeed()[0]
	fakePieces, _ := strconv.Atoi(fake.DigitsN(2))
	fakeDate := time.Date(fake.Year(2022, 2023), time.Month(fake.MonthNum()), fake.Day(), 0, 0, 0, 0, time.UTC)
	return models.PharmacyStock{
		MedicineKey: medicine.Key,
		LotNumber:   fake.DigitsN(9),
		Pieces:      int16(fakePieces),
		ExpiresAt:   fakeDate,
	}
}
