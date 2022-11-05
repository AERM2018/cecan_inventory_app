package models

import (
	"time"

	"github.com/google/uuid"
)

type FixedAssetDescription struct {
	Id          uuid.UUID `gorm:"primaryKey,default:'uuid_generate_v4()'" json:"id"`
	Description string    `json:"description"`
	Model       string    `json:"model"`
	Brand       string    `json:"brand"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
