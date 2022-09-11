package mocks

import (
	"cecan_inventory/domain/models"
	"strconv"
	"time"

	"github.com/icrowley/fake"
)

func GetRolesMock() []models.Role {
	return []models.Role{
		{Name: "Medico"},
		{Name: "Farmacia"},
		{Name: "Almacén"},
		{Name: "Estandar"},
		{Name: "Ingeniero"},
		{Name: "Auditor"},
		{Name: "Admin"},
	}
}

func GetUserMock(rolId string) models.User {
	return models.User{
		RoleId:   rolId,
		Password: "Qwerty*123",
		Name:     "CECAN ADMIN",
		Email:    "admin@cecan.com",
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
