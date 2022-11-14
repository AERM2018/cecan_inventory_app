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
		DepartmentId uuid.UUID                  `json:"department_id"`
		FixedAssets  []FixedAssetsItemsRequests `json:"fixed_assets"`
		UserId       string                     `json:"user_id"`
		CreatedAt    time.Time                  `json:"created_at"`
		UpdatedAt    time.Time                  `json:"updated_at"`
	}
)
