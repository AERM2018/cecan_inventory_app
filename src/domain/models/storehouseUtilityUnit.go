package models

import "github.com/google/uuid"

type StorehouseUtilityUnit struct {
	Id   uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name string    `json:"name"`
}
