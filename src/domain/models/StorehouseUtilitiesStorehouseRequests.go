package models

import (
	"time"

	"github.com/google/uuid"
)

type StorehouseUtilitiesStorehouseRequests struct {
	Id                  uuid.UUID                 `gorm:"primaryKey;default:'uuid_generate_v4();'" json:"id"`
	UtilityKey          string                    `gorm:"column:storehouse_utility_key" json:"key"`
	StorehouseUtilty    StorehouseUtilityDetailed `gorm:"foreignKey:storehouse_utility_key" json:"details,omitempty"`
	StorehouseRequestId uuid.UUID                 `gorm:"column:storehouse_request_id;foreignKey" json:"request_id"`
	Pieces              int16                     `json:"pieces"`
	PiecesSupplied      int16                     `json:"pieces_supplied"`
	CreatedAt           *time.Time                `json:"created_at,omitempty"`
	UpdatedAt           *time.Time                `json:"updated_at,omitempty"`
}
