package mocks

import (
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	authtoken "cecan_inventory/infrastructure/external/authToken"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/icrowley/fake"
)

func GetTokenMock(claims models.AuthClaims) string {
	token, _ := authtoken.GenerateJWT(claims)
	return token
}

func GetRolesMock(name string) []models.Role {
	rolesMocksMap := []map[string]string{
		{"id": "7d3f3faa-39e2-4b3c-aaa8-8ca60fa090b4", "name": "Medico"},
		{"id": "6648c302-f487-491d-be57-5f90bbe380c7", "name": "Farmacia"},
		{"id": "5b4e7720-b663-474b-841d-aed481907908", "name": "Almacén"},
		{"id": "6542908e-41ca-4230-8b45-985be3392b5d", "name": "Estandar"},
		{"id": "097305e0-c4a3-4fd6-a929-49cd1aca5e90", "name": "Ingeniero"},
		{"id": "e2660b8f-404c-4baa-b7f5-34e98c602046", "name": "Auditor"},
		{"id": "3c9c3b1c-80c5-43a0-9276-4f88b67a7ac7", "name": "Admin"},
	}
	var rolesMocks []models.Role
	var rolesFinalMap []map[string]string
	if name == "" {
		rolesFinalMap = rolesMocksMap
	} else {
		rolesFound := common.FilterSlice(rolesMocksMap, func(i interface{}) bool {
			parsed := i.(map[string]string)
			return strings.ToLower(parsed["name"]) == strings.ToLower(name)
		})
		rolesFinalMap = rolesFound.([]map[string]string)
	}
	for _, roleMockAsMap := range rolesFinalMap {
		idAsUuid, _ := uuid.Parse(roleMockAsMap["id"])
		rolesMocks = append(rolesMocks, models.Role{Id: &idAsUuid, Name: roleMockAsMap["name"]})
	}
	return rolesMocks
}

func GetUserMockSeed(rolId string) models.User {
	return models.User{
		RoleId:   rolId,
		Password: "Qwerty*123",
		Name:     "CECAN ADMIN",
		Surname:  "super",
		Email:    "admin@cecan.com",
	}
}

func GetUserMock(rolId string, minPassLen int) models.User {
	fakeEmail, _, _ := strings.Cut(fake.EmailAddress(), "@")
	return models.User{
		RoleId:   rolId,
		Password: fake.Password(minPassLen, minPassLen, true, true, false),
		Name:     fake.MaleFirstName(),
		Surname:  fake.MaleLastName(),
		Email:    fmt.Sprintf("%v@cecan.com", fakeEmail),
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

func GetPharmacyStockMock(medicine models.Medicine) models.PharmacyStock {
	if (medicine == models.Medicine{}) {
		medicine = GetMedicineMockSeed()[0]
	}
	fakePieces, _ := strconv.Atoi(fake.DigitsN(2))
	fakeDate := time.Date(fake.Year(2023, 2024), time.Month(fake.MonthNum()), fake.Day(), 0, 0, 0, 0, time.UTC)
	return models.PharmacyStock{
		Id:          uuid.New(),
		MedicineKey: medicine.Key,
		LotNumber:   fake.DigitsN(9),
		Pieces:      int16(fakePieces),
		ExpiresAt:   fakeDate,
	}
}

func GetPharmacyStockMockSeed() []models.PharmacyStock {
	pointer := 0
	fakeUuids := []string{"15af71ce-1183-498c-a744-bbbc4d181010", "1646b651-4545-4f88-851d-f61e780c8d8a"}
	pharmacyStocksMocksSeed := make([]models.PharmacyStock, 0)
	for pointer < 2 {
		fakePieces, _ := strconv.Atoi(fake.DigitsN(2))
		fakeDate := time.Date(fake.Year(2023, 2024), time.Month(fake.MonthNum()), fake.Day(), 0, 0, 0, 0, time.UTC)
		uuidParsed, _ := uuid.Parse(fakeUuids[pointer])
		pharmacyStockMock := models.PharmacyStock{
			Id:          uuidParsed,
			MedicineKey: GetMedicineMockSeed()[0].Key,
			LotNumber:   fake.DigitsN(9),
			Pieces:      int16(fakePieces),
			ExpiresAt:   fakeDate,
		}
		pharmacyStocksMocksSeed = append(pharmacyStocksMocksSeed, pharmacyStockMock)
		pointer += 1
	}
	// Changes pieces used to the last pharmacy stock for testing pursposes
	pharmacyStocksMocksSeed[len(pharmacyStocksMocksSeed)-1].Pieces_used = 2
	pharmacyStocksMocksSeed[len(pharmacyStocksMocksSeed)-1].Pieces -= pharmacyStocksMocksSeed[len(pharmacyStocksMocksSeed)-1].Pieces_used
	return pharmacyStocksMocksSeed
}

func GetPrescriptionMockSeed() []models.PrescriptionDetialed {
	pointer := 0
	fakeUuids := []string{"237aa448-3e83-4af0-ae24-e1b1138f6fec", "f64cc4d4-9c33-4a97-a981-7539b74fc07b"}
	prescriptionsMocksSeed := make([]models.PrescriptionDetialed, 0)
	for pointer < 2 {
		fakePieces, _ := strconv.Atoi(fake.DigitsN(1))
		uuidParsed, _ := uuid.Parse(fakeUuids[pointer])
		prescriptionMock := models.PrescriptionDetialed{
			Id:           uuidParsed,
			UserId:       GetUserMock(GetRolesMock("medico")[0].Id.String(), 10).Id,
			PatientName:  fake.FullName(),
			Instructions: fake.Paragraph(),
			Medicines: []models.PrescriptionsMedicines{
				{
					MedicineKey: GetMedicineMockSeed()[0].Key,
					Pieces:      int16(fakePieces),
				},
			},
		}
		prescriptionsMocksSeed = append(prescriptionsMocksSeed, prescriptionMock)
		pointer += 1
	}
	return prescriptionsMocksSeed
}
