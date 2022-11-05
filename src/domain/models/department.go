package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Department struct {
		Id                uuid.UUID `gorm:"primaryKey,defualt:'uuid_generate_v4()'" json:"id"`
		ResponsibleUserId *string   `json:"resposible_user_id,omitempty"`
		Name              string    `json:"name"`
		FloorNumber       string    `json:"floor_number"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	}

	DepartmentDetailed struct {
		Id                uuid.UUID `gorm:"primaryKey,defualt:'uuid_generate_v4()'" json:"id"`
		ResponsibleUserId string    `json:"resposible_user_id"`
		ResponsibleUser   User      `gorm:"foreignKey" json:"reponsible_user"`
		Name              string    `json:"name"`
		FloorNumber       string    `json:"floor_number"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	}
)
