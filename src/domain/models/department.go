package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Department struct {
		Id                uuid.UUID      `gorm:"primaryKey;default:'uuid_generate_v4()'" json:"id"`
		ResponsibleUserId *string        `json:"resposible_user_id"`
		Name              string         `json:"name"`
		FloorNumber       string         `json:"floor_number"`
		CreatedAt         *time.Time     `json:"created_at,omitempty"`
		UpdatedAt         *time.Time     `json:"updated_at,omitempty"`
		DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	}

	DepartmentDetailed struct {
		Id                uuid.UUID      `json:"id"`
		ResponsibleUserId string         `json:"resposible_user_id"`
		ResponsibleUser   User           `gorm:"foreignKey:ResponsibleUserId" json:"reponsible_user"`
		Name              string         `json:"name"`
		FloorNumber       string         `json:"floor_number"`
		CreatedAt         time.Time      `json:"created_at,omitempty"`
		UpdatedAt         time.Time      `json:"updated_at,omitempty"`
		DeletedAt         gorm.DeletedAt `json:"deleted_at"`
	}
)

func (department *Department) BeforeCreate(tx *gorm.DB) error {
	department.Name = strings.ToUpper(department.Name)
	department.FloorNumber = strings.ToUpper(department.FloorNumber)
	return nil
}

func (department *Department) BeforeUpdate(tx *gorm.DB) error {
	department.Name = strings.ToUpper(department.Name)
	department.FloorNumber = strings.ToUpper(department.FloorNumber)
	return nil
}
