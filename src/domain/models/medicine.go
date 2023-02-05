package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Medicine struct {
		Key       string         `gorm:"primaryKey" json:"key" validate:"required"`
		Name      string         `json:"name" validate:"required"`
		CreatedAt *time.Time     `json:"created_at,omitempty"`
		UpdatedAt *time.Time     `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	}

	
	MedicinesFilters struct {
		MedicineKey string `json:"medicine_key,omitempty"`
		MedicineName string `json:"medicine_name,omitempty"`
		Limit int `json:"limit,omitempty"`
		Page int `json:"page,omitempty"`
		IncludeDeleted bool `json:"include_deleted,omitempty"`
		ShowLessQty bool `json:"show_less_qty,omitempty"`
	}
)
