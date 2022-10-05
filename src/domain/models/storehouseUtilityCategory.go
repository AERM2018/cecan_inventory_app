package models

import "github.com/google/uuid"

type StorehouseUtilityCategory struct {
	Id   uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name string    `json:"name"`
}
