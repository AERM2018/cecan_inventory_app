package models

import "github.com/google/uuid"

type Role struct {
	Id   *uuid.UUID `gorm:"primaryKey;default:'uuid_generate_v4()'" json:"id,omitempty"`
	Name string    `json:"name"`
}
