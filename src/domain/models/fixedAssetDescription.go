package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FixedAssetDescription struct {
	Id          uuid.UUID `gorm:"primaryKey;default:'uuid_generate_v4();'" json:"id"`
	Description string    `json:"description"`
	Model       string    `json:"model"`
	Brand       string    `json:"brand"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (fixedAssetDescription *FixedAssetDescription) BeforeCreate(tx *gorm.DB) error {
	fixedAssetDescription.Description = strings.ToUpper(fixedAssetDescription.Description)
	fixedAssetDescription.Brand = strings.ToUpper(fixedAssetDescription.Brand)
	fixedAssetDescription.Model = strings.ToUpper(fixedAssetDescription.Model)
	return nil
}
