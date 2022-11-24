package datasources

import (
	"cecan_inventory/domain/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorehouseRequestsDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc StorehouseRequestsDataSource) GetStorehouseRequests() ([]models.StorehouseRequestDetailed, error) {
	storehouseRequests := make([]models.StorehouseRequestDetailed, 0)
	err := dataSrc.DbPsql.
		Table("storehouse_requests").
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "name", "surname", "role_id")
		}).
		Preload("User.Role").
		Preload("Status", func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("created_at", "updated_at")
		}).
		Find(&storehouseRequests).Error
	if err != nil {
		return storehouseRequests, err
	}
	return storehouseRequests, nil
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
			return tx.Table("storehouse_utilities").Select("key", "generic_name", "storehouse_utility_category_id", "storehouse_utility_unit_id", "storehouse_utility_presentation_id")
		}).
		Preload("Utilities.StorehouseUtilty.Category", func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("created_at", "updated_at")
		}).
		Preload("Utilities.StorehouseUtilty.Unit", func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("created_at", "updated_at")
		}).
		Preload("Utilities.StorehouseUtilty.Presentation", func(tx *gorm.DB) *gorm.DB {
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

func (dataSrc StorehouseRequestsDataSource) PutUtilityIntoRequest(storehouseRequestUtilty models.StorehouseUtilitiesStorehouseRequests) error {
	err := dataSrc.DbPsql.
		Create(&storehouseRequestUtilty).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc StorehouseRequestsDataSource) UpdateStorehouseRequest(id string, requestInfo models.StorehouseRequest, utilitiesOperationTypes ...[]models.StorehouseUtilitiesStorehouseRequests) (uuid.UUID, error) {
	if len(utilitiesOperationTypes) != 3 {
		panic("The utilities slice must have a length of three: utilities to update, to add and to remove")
	}
	errInTransaction := dataSrc.DbPsql.Transaction(func(tx *gorm.DB) error {
		var err error
		// Update storehouse request instance with basic information
		err = tx.
			Table("storehouse_requests").
			Omit("id", "created_at").
			Where("id = ?", id).
			Updates(&requestInfo).
			Error
		if err != nil {
			fmt.Println(err)
			return errors.New("No se pudo actualizae la solicitud, verifique los datos y vuelvalo a intentar.")
		}
		// Start adding or removing utilities to the storehouse request
		// utilitiesToAdd, utilitiesToRemove := requestDetailed.FilterUtilitesFromPrescription()
		// The slice's position 0 is for the utilities to be updated that already exists in the storehouse request
		for _, storehouseUtilitiesFromRequest := range utilitiesOperationTypes[0] {
			err = tx.
				Table("storehouse_utilities_storehouse_requests").
				Where("storehouse_request_id = ? AND storehouse_utility_key = ?", id, storehouseUtilitiesFromRequest.UtilityKey).
				Updates(&storehouseUtilitiesFromRequest).Error
			if err != nil {
				return errors.New("No se pudo actualizar la solicitud debido a que no se pudo actualizar los elementos del almacen que le pertencen a la misma.")
			}
		}
		// The slice's position 1 is for the utilities to be added to the storehouse request
		for _, utilityToAdd := range utilitiesOperationTypes[1] {
			err = tx.
				Create(&utilityToAdd).Error
			if err != nil {
				return errors.New("No se pudo añadir un nuevo elemento de almacen a la solictud al actualizarla.")
			}
		}
		// The slice's position 2 is for the utilities to be removed to the storehouse request
		for _, utilityToRemove := range utilitiesOperationTypes[2] {
			err = tx.
				Delete(&utilityToRemove).Error
			if err != nil {
				return errors.New("No se pudo remover un elemento de almacen de la solicitud al actualizarla.")
			}
		}
		return nil
	})
	if errInTransaction != nil {
		return uuid.Nil, errInTransaction
	}
	return requestInfo.Id, nil
}

func (dataSrc StorehouseRequestsDataSource) DeleteStorehouseRequest(id string) error {
	errInTransaction := dataSrc.DbPsql.Transaction(func(tx *gorm.DB) error {

		errDeletingUtilities := tx.
			Where("storehouse_request_id = ?", id).
			Delete(&models.StorehouseUtilitiesStorehouseRequests{}).
			Error
		if errDeletingUtilities != nil {
			return errors.New("No se pudo desasignar los elementos de almacen de la solicitud.")
		}
		errDeletingRequest := tx.
			Where("id = ?", id).
			Delete(&models.StorehouseRequest{}).
			Error
		if errDeletingRequest != nil {
			return errors.New("No se pudo la solicitud, intentelo de nuevo.")
		}
		return nil
	})
	if errInTransaction != nil {
		return errInTransaction
	}
	return nil
}

func (dataSrc StorehouseRequestsDataSource) IsSameRequestCreator(id string, userId string) bool {
	var storehouseRequest models.StorehouseRequest
	dataSrc.DbPsql.
		Where("id = ?", id).
		Take(&storehouseRequest)
	return strings.ToLower(storehouseRequest.UserId) == strings.ToLower(userId)
}

func (dataSrc StorehouseRequestsDataSource) IsRequestDeterminedStatus(id string, statuses []string) bool {
	var storehouseRequest models.StorehouseRequestDetailed
	dataSrc.DbPsql.
		Table("storehouse_requests").
		Joins("Status", dataSrc.DbPsql.Where("name in (?)", statuses).Model(&models.StorehouseRequestStatus{})).
		Where("storehouse_requests.id = ?", id).
		Take(&storehouseRequest)
	return uuid.Nil != storehouseRequest.Status.Id
}

func (dataSrc StorehouseRequestsDataSource) SupplyStorehouseRequest(id string, utilities []models.StorehouseUtilitiesStorehouseRequests) error {
	errInTransaction := dataSrc.DbPsql.Transaction(func(tx *gorm.DB) error {
		var (
			utilityRequested        models.StorehouseUtilitiesStorehouseRequests
			storehouseRequestStatus models.StorehouseRequestStatus
			isRequestIncomplete     bool
		)

		for _, utilityRequestReference := range utilities {
			errConsulting := tx.
				Table("storehouse_utilities_storehouse_requests").
				Where("storehouse_request_id = ? and storehouse_utility_key = ?", id, utilityRequestReference.UtilityKey).
				Take(&utilityRequested).Error
			if errConsulting != nil {
				return errConsulting
			}
			if utilityRequested.PiecesSupplied+utilityRequestReference.PiecesSupplied < utilityRequested.Pieces {
				isRequestIncomplete = true
			}
			errUpdating := tx.
				Table("storehouse_utilities_storehouse_requests").
				Where("storehouse_request_id = ? and storehouse_utility_key = ?", id, utilityRequestReference.UtilityKey).
				Update("pieces_supplied", utilityRequested.PiecesSupplied+utilityRequestReference.PiecesSupplied).
				Update("last_pieces_supplied", utilityRequestReference.PiecesSupplied).
				Error
			if errUpdating != nil {
				return errUpdating
			}
			errInUtilitiesResetvation := tx.Transaction(func(tx *gorm.DB) error {
				// Take and reserver the utilities stocks for the request
				err := tx.Exec(fmt.Sprintf("Call public.reserve_utility_to_request('%v');", id)).Error
				if err != nil {
					return err
				}
				return nil
			})
			if errInUtilitiesResetvation != nil {
				return errInUtilitiesResetvation
			}
		}
		if isRequestIncomplete {
			tx.Where("name = ?", "Incompleta").Take(&storehouseRequestStatus)
		} else {
			tx.Where("name = ?", "Completada").Take(&storehouseRequestStatus)
		}
		errUpdatingStatus := tx.
			Table("storehouse_requests").
			Where("id = ?", id).
			Update("storehouse_request_status_id", storehouseRequestStatus.Id).
			Update("supplied_at", time.Now()).
			Error
		if errUpdatingStatus != nil {
			return errors.New("No se pudo actualizar el estado de la solictud de almacen.")
		}
		return nil
	})
	if errInTransaction != nil {
		errorMsg := ""
		_, errorMsg, _ = strings.Cut(errInTransaction.Error(), ":")
		errorMsg, _, _ = strings.Cut(errorMsg, "(SQLSTATE")
		errorMsg = strings.Trim(errorMsg, " ")
		return errors.New(errorMsg)
	}
	return nil
}
