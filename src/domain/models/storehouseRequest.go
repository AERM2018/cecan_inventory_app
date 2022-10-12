package models

import (
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
		Utilities []StorehouseUtilitiesStorehouseRequests `gorm:"foreignKey:storehouse_request_id" json:"utilities,omitempty"`
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
