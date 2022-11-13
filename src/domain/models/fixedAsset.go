package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	FixedAsset struct {
		Key                     string     `gorm:"primaryKey" json:"key"`
		Description             string     `json:"description"`
		Brand                   string     `json:"brand"`
		Model                   string     `json:"model"`
		FixedAssetDescriptionId *uuid.UUID `json:"description_id,omitempty"`
		Series                  string     `json:"series"`
		Type                    string     `json:"type"`
		PhysicState             string     `json:"physic_state"`
		DepartmentId            string     `json:"department_id"`
		Observation             string     `json:"observations"`
		DirectorUserId          string     `json:"director_user_id"`
		AdministratorUserId     string     `json:"administrator_user_id"`
		CreatedAt               time.Time  `json:"created_at"`
		UpdatedAt               time.Time  `json:"updated_at"`
	}

	FixedAssetDetailed struct {
		Key                   string    `gorm:"primaryKey" json:"key"`
		Description           string    `json:"description"`
		Brand                 string    `json:"brand"`
		Model                 string    `json:"model"`
		Series                string    `json:"series"`
		Type                  string    `json:"type"`
		PhysicState           string    `json:"physic_state"`
		DepartmentId          string    `json:"department_id"`
		DepartmentName        string    `json:"department_name"`
		Observation           string    `json:"observations"`
		DirectorUserId        string    `json:"director_user_id"`
		DirectorUserName      string    `json:"director_user_name"`
		AdministratorUserId   string    `json:"administrator_user_id"`
		AdministratorUserName string    `json:"administrator_user_name"`
		CreatedAt             time.Time `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
	}

	FixedAssetFilters struct {
		Brand          string `json:"brand,omitempty"`
		Model          string `json:"model,omitempty"`
		Type           string `json:"type,omitempty"`
		PhysicState    string `json:"physic_state,omitempty"`
		DepartmentName string `json:"department_name,omitempty"`
	}
)
