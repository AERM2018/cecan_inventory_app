package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	FixedAssetsRequest struct {
		Id           uuid.UUID `gorm:"primaryKey;default:'uuid_generate_v4();'" json:"id"`
		DepartmentId uuid.UUID `json:"department_id"`
		UserId       string    `json:"user_id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	FixedAssetsRequestDetailed struct {
		Id           uuid.UUID                  `gorm:"primaryKey;default:'uuid_generate_v4();'" json:"id"`
		DepartmentId uuid.UUID                  `gorm:"foreignKey:department_id" json:"department_id"`
		Department   DepartmentDetailed         `json:"department"`
		FixedAssets  []FixedAssetsItemsRequests `gorm:"foreignKey:fixed_assets_request_id" json:"fixed_assets"`
		UserId       string                     `jorm:"foreignKey:user_id" json:"user_id"`
		User         User                       `json:"user"`
		CreatedAt    time.Time                  `json:"created_at"`
		UpdatedAt    time.Time                  `json:"updated_at"`
	}
)
