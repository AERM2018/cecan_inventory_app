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
	pharmacyStocksMocksSeed[0].ExpiresAt = time.Now().UTC().Add(time.Hour * 24)
	// Changes pieces used to the last pharmacy stock for testing pursposes
	pharmacyStocksMocksSeed[len(pharmacyStocksMocksSeed)-1].Pieces_used = 2
	pharmacyStocksMocksSeed[len(pharmacyStocksMocksSeed)-1].Pieces -= pharmacyStocksMocksSeed[len(pharmacyStocksMocksSeed)-1].Pieces_used
	return pharmacyStocksMocksSeed
}

func GetPrescriptionStatuesMockSeed() []models.PrescriptionsStatues {
	prescriptionStatues := make([]models.PrescriptionsStatues, 0)
	mapStatues := []map[string]string{
		{
			"id":   "6f328cc7-3662-4d6f-912a-34449c241019",
			"name": "Pendiente",
		},
		{
			"id":   "8036632b-d5aa-4cf7-81ed-b4abbbd90482",
			"name": "Completada",
		},
	}
	for _, prescriptionStatus := range mapStatues {
		uuidParsed, _ := uuid.Parse(prescriptionStatus["id"])
		status := models.PrescriptionsStatues{
			Id:   uuidParsed,
			Name: prescriptionStatus["name"],
		}
		prescriptionStatues = append(prescriptionStatues, status)
	}
	return prescriptionStatues
}
func GetPrescriptionMockSeed() []models.PrescriptionDetialed {
	pointer := 0
	fakeUuids := []string{"237aa448-3e83-4af0-ae24-e1b1138f6fec", "f64cc4d4-9c33-4a97-a981-7539b74fc07b", "5db2f10a-d692-483b-8ba6-3ea48d2f00c6"}
	prescriptionsMocksSeed := make([]models.PrescriptionDetialed, 0)
	for pointer < 3 {
		fakePieces, _ := strconv.Atoi(fake.DigitsN(1))
		uuidParsed, _ := uuid.Parse(fakeUuids[pointer])
		prescriptionMock := models.PrescriptionDetialed{
			Id:           uuidParsed,
			UserId:       GetUserMock(GetRolesMock("medico")[0].Id.String(), 10).Id,
			PatientName:  fake.FullName(),
			Instructions: fake.Paragraph(),
			Medicines: []models.PrescriptionsMedicines{
				{
					PrescriptionId: uuidParsed,
					MedicineKey:    GetMedicineMockSeed()[0].Key,
					Pieces:         int16(fakePieces),
				},
			},
		}
		prescriptionsMocksSeed = append(prescriptionsMocksSeed, prescriptionMock)
		pointer += 1
	}
	// Set the last prescription as completed
	isStatus, status := common.FindInSlice(GetPrescriptionStatuesMockSeed(), func(i interface{}) bool {
		parsed := i.(models.PrescriptionsStatues)
		return strings.ToLower(parsed.Name) == "completada"
	})
	if isStatus {
		prescriptionsMocksSeed[len(prescriptionsMocksSeed)-1].PrescriptionStatusId = status.([]models.PrescriptionsStatues)[0].Id
	}
	return prescriptionsMocksSeed
}

func GetPrescriptionMock() models.PrescriptionDetialed {
	fakePieces, _ := strconv.Atoi(fake.DigitsN(1))
	uuidParsed := uuid.New()
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
	return prescriptionMock
}

func GetStorehouseUtiltyCategoriesMockSeed() []models.StorehouseUtilityCategory {
	pointer := 0
	mapStorehouseUtilityCategories := []map[string]string{
		{"id": "ff0b0d3a-d013-4946-9cc2-285a16861ec1", "name": "Material medico"},
		{"id": "52fc404c-dd25-4d51-9e93-adbc4329f6d0", "name": "Material de limpieza"},
	}
	storehouseUtilityCategoriesMockSeed := make([]models.StorehouseUtilityCategory, 0)
	for pointer < len(mapStorehouseUtilityCategories) {
		for _, storehouseUtilityCategory := range mapStorehouseUtilityCategories {
			uuidParsed, _ := uuid.Parse(storehouseUtilityCategory["id"])
			status := models.StorehouseUtilityCategory{
				Id:   uuidParsed,
				Name: storehouseUtilityCategory["name"],
			}
			storehouseUtilityCategoriesMockSeed = append(storehouseUtilityCategoriesMockSeed, status)
		}
		pointer += 1
	}
	return storehouseUtilityCategoriesMockSeed
}

func GetStorehouseUtiltyPresentationsMockSeed() []models.StorehouseUtilityPresentation {
	pointer := 0
	mapStorehouseUtilityPresentations := []map[string]string{
		{"id": "88f6e5d7-f836-4079-bfb0-9e4199a248e2", "name": "Botella"},
		{"id": "0b17934b-967a-4f2f-92dc-a0632f1f19ed", "name": "Caja"},
	}
	storehouseUtilityPresentationsMockSeed := make([]models.StorehouseUtilityPresentation, 0)
	for pointer < len(mapStorehouseUtilityPresentations) {
		for _, storehouseUtilityPresentation := range mapStorehouseUtilityPresentations {
			uuidParsed, _ := uuid.Parse(storehouseUtilityPresentation["id"])
			presentation := models.StorehouseUtilityPresentation{
				Id:   uuidParsed,
				Name: storehouseUtilityPresentation["name"],
			}
			storehouseUtilityPresentationsMockSeed = append(storehouseUtilityPresentationsMockSeed, presentation)
		}
		pointer += 1
	}
	return storehouseUtilityPresentationsMockSeed
}

func GetStorehouseUtiltyUnitsMockSeed() []models.StorehouseUtilityUnit {
	pointer := 0
	mapStorehouseUtilityUnits := []map[string]string{
		{"id": "f528523e-89ac-41e5-8567-ede21ab37f13", "name": "Litro(s)"},
		{"id": "a569159b-e4af-4072-ab65-a4e39d8b0f07", "name": "Pieza(s)"},
	}
	storehouseUtilityUnitsMockSeed := make([]models.StorehouseUtilityUnit, 0)
	for pointer < len(mapStorehouseUtilityUnits) {
		for _, storehouseUtilityUnit := range mapStorehouseUtilityUnits {
			uuidParsed, _ := uuid.Parse(storehouseUtilityUnit["id"])
			presentation := models.StorehouseUtilityUnit{
				Id:   uuidParsed,
				Name: storehouseUtilityUnit["name"],
			}
			storehouseUtilityUnitsMockSeed = append(storehouseUtilityUnitsMockSeed, presentation)
		}
		pointer += 1
	}
	return storehouseUtilityUnitsMockSeed
}

func GetStorehouseRequestStatuses() []models.StorehouseRequestStatus {
	pointer := 0
	mapStorehouseRequestStatuses := []map[string]string{
		{"id": "8d60d759-55e3-474e-9b13-bd41c720ef2b", "name": "Pendiente"},
		{"id": "4920d64d-7c5c-4fb0-ad52-c00141a298b9", "name": "Completada"},
	}
	storehouseRequestStatusesMockSeed := make([]models.StorehouseRequestStatus, 0)
	for pointer < len(mapStorehouseRequestStatuses) {
		newUuid, _ := uuid.Parse(mapStorehouseRequestStatuses[pointer]["id"])
		storehouseRequestStatusesMockSeed = append(
			storehouseRequestStatusesMockSeed,
			models.StorehouseRequestStatus{
				Id:   newUuid,
				Name: mapStorehouseRequestStatuses[pointer]["name"],
			})
		pointer += 1
	}
	return storehouseRequestStatusesMockSeed
}
