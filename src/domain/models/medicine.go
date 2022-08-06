package models

import (
	"time"

	"gorm.io/gorm"
)

type Medicine struct {
	Key       string         `gorm:"primaryKey" json:"key" validate:"required"`
	Name      string         `json:"name" validate:"required"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
