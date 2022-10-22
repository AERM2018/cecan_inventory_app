package models

import (
	"cecan_inventory/domain/common"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	StorehouseRequest struct {
		Id        uuid.UUID  `gorm:"primaryKey;default:'uuid_generate_v4();'" json:"id"`
		UserId    string     `json:"user_id"`
		Folio     int16      `json:"folio,omitempty"`
		StatusId  uuid.UUID  `gorm:"column:storehouse_request_status_id" json:"status_id"`
		CreatedAt *time.Time `json:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
	}

	StorehouseRequestDetailed struct {
		Id        uuid.UUID                               `gorm:"primaryKey;default:'uuid_generate_v4();'" json:"id"`
		UserId    string                                  `json:"user_id"`
		User      User                                    `gorm:"foreignKey:user_id" json:"user,omitempty"`
		Folio     int16                                   `json:"folio"`
		StatusId  uuid.UUID                               `gorm:"column:storehouse_request_status_id" json:"status_id"`
		Status    StorehouseRequestStatus                 `json:"status,omitempty"`
		Utilities []StorehouseUtilitiesStorehouseRequests `gorm:"foreignKey:storehouse_request_id" json:"utilities"`
		CreatedAt *time.Time                              `json:"created_at,omitempty"`
		UpdatedAt *time.Time                              `json:"updated_at,omitempty"`
	}
)

func (request *StorehouseRequest) BeforeCreate(tx *gorm.DB) error {
	var (
		lastRequest StorehouseRequest
		status      StorehouseRequestStatus
		err         error
	)
	if request.StatusId == uuid.Nil {
		err = tx.
			Where("name = ?", "Pendiente").
			Take(&status).Error
		if err != nil {
			fmt.Printf("Error before create: %v", err)
			return err
		}
		request.StatusId = status.Id
	}
	err = tx.
		Model(&StorehouseRequest{}).
		Order("folio DESC").
		First(&lastRequest).Error

	request.Folio = lastRequest.Folio + 1
	return nil
}

func (storehouseRequestDetailed *StorehouseRequestDetailed) FilterUtilitesFromRequest(oldUtilites []StorehouseUtilitiesStorehouseRequests) (
	utilitiesToAdd []StorehouseUtilitiesStorehouseRequests,
	utilitiesToRemove []StorehouseUtilitiesStorehouseRequests,
) {
	filteredStorehouseUtilities := make([]StorehouseUtilitiesStorehouseRequests, 0)
	storehouseUtilitiesToInsert := make([]StorehouseUtilitiesStorehouseRequests, 0)
	storehouseUtilitiesToDelete := make([]StorehouseUtilitiesStorehouseRequests, 0)
	for _, newStorehouseUtilityRequest := range storehouseRequestDetailed.Utilities {
		isUtility, storehouseUtilityRequest := common.FindInSlice(oldUtilites, func(i interface{}) bool {
			parsed := i.(StorehouseUtilitiesStorehouseRequests)
			return newStorehouseUtilityRequest.UtilityKey == parsed.UtilityKey

		})
		// Insert utilty if it's not comming with the utilities already assigned to the request
		if !isUtility {
			newStorehouseUtilityRequest.StorehouseRequestId = storehouseRequestDetailed.Id
			storehouseUtilitiesToInsert = append(storehouseUtilitiesToInsert, newStorehouseUtilityRequest)
		}
		// Update pieces requested of an utility which's already assigned
		if isUtility && newStorehouseUtilityRequest.Pieces != storehouseUtilityRequest.([]StorehouseUtilitiesStorehouseRequests)[0].Pieces {
			filteredStorehouseUtilities = append(filteredStorehouseUtilities, newStorehouseUtilityRequest)
		}
	}

	for _, oldStorehouseUtilityRequest := range oldUtilites {
		isUtility, _ := common.FindInSlice(storehouseRequestDetailed.Utilities, func(i interface{}) bool {
			parsed := i.(StorehouseUtilitiesStorehouseRequests)
			return oldStorehouseUtilityRequest.UtilityKey == parsed.UtilityKey

		})
		// Remove utility if the update utility list does not include a utility already assigned to the request
		if !isUtility {
			storehouseUtilitiesToDelete = append(storehouseUtilitiesToDelete, oldStorehouseUtilityRequest)
		}
	}
	storehouseRequestDetailed.Utilities = filteredStorehouseUtilities
	return storehouseUtilitiesToInsert, storehouseUtilitiesToDelete
}
