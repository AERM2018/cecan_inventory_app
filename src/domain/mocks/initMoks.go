package mocks

import (
	"cecan_inventory/domain/models"

	"github.com/icrowley/fake"
)

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
