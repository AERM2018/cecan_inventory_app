package datasources

import (
	"cecan_inventory/domain/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorehouseRequestsDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc StorehouseRequestsDataSource) GetStorehouseRequestById(id string) (models.StorehouseRequestDetailed, error) {
	storehouseRequest := models.StorehouseRequestDetailed{}
	err := dataSrc.DbPsql.
		Table("storehouse_requests").
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "name", "surname", "role_id")
		}).
		Preload("User.Role").
		Preload("Status", func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("created_at", "updated_at")
		}).
		Preload("Utilities", func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("created_at", "updated_at")
		}).
		Preload("Utilities.StorehouseUtilty", func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("created_at", "updated_at")
		}).
		Where("id = ?", id).
		Take(&storehouseRequest).Error
	if err != nil {
		return storehouseRequest, err
	}
	return storehouseRequest, nil
}

func (dataSrc StorehouseRequestsDataSource) CreateStorehouseRequest(requestInfo models.StorehouseRequest, requestUtilities []models.StorehouseUtilitiesStorehouseRequests) (uuid.UUID, error) {
	errInTransaction := dataSrc.DbPsql.Transaction(func(tx *gorm.DB) error {
		var err error
		// Create a new storehouse request instance with basic information
		err1 := tx.
			Create(&requestInfo).
			Error
		if err1 != nil {
			fmt.Println(err1)
			return errors.New("No se pudo crear una nueva solicitud, verifique los datos y vuelvalo a intentar.")
		}
		// Start assosiationg utilities to the new storehouse request
		for _, storehouseUtilitiesFromRequest := range requestUtilities {
			storehouseUtilitiesFromRequest.StorehouseRequestId = requestInfo.Id
			err = tx.Create(&storehouseUtilitiesFromRequest).Error
			if err != nil {
				return errors.New("No se pudo crear la solicitud debido a que no se pudo asignar los elementos del almacen a la misma.")
			}
		}
		return nil
	})
	if errInTransaction != nil {
		return uuid.Nil, errInTransaction
	}
	return requestInfo.Id, nil
}
