package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrescriptionsStatues struct {
	Id        uuid.UUID      `gorm:"primaryKey" json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}
