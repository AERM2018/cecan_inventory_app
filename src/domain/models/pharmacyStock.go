package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SemaforizationColor string

const (
	Red   SemaforizationColor = "red"
	Ambar SemaforizationColor = "ambar"
	Green SemaforizationColor = "green"
)

type PharmacyStock struct {
	Id                  uuid.UUID           `gorm:"primaryKey;default:'uuid_generate_v4()'" json:"id"`
	MedicineKey         string              `gorm:"foreignKey:medicine_key;references:key" json:"medicine_key" validate:"required"`
	Medicine            *Medicine           `gorm:"foreignKey:MedicineKey" json:"medicine,omitempty"`
	LotNumber           string              `json:"lot_number" validate:"required"`
	Pieces              int16               `json:"pieces" validate:"required"`
	SemaforizationColor SemaforizationColor `json:"semaforization_color"`
	CreatedAt           *time.Time          `json:"created_at,omitempty"`
	UpdatedAt           *time.Time          `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`
	ExpiresAt           time.Time           `json:"expires_at" validate:"required"`
	DeletedAt           gorm.DeletedAt      `gorm:"index" json:"deleted_at"`
}

func (pharmacyStock *PharmacyStock) BeforeCreate(tx *gorm.DB) (err error) {
	if pharmacyStock.ExpiresAt.Unix() < time.Now().AddDate(0, 6, 0).Unix() {
		pharmacyStock.SemaforizationColor = Red
	} else if pharmacyStock.ExpiresAt.Unix() >= time.Now().AddDate(0, 6, 0).Unix() && pharmacyStock.ExpiresAt.Unix() <= time.Now().AddDate(0, 12, 0).Unix() {
		pharmacyStock.SemaforizationColor = Ambar
	} else {
		pharmacyStock.SemaforizationColor = Green
	}
	return
}

type PharmacyStocksDetailed struct {
	Medicine                    Medicine        `json:"medicine"`
	Stocks                      []PharmacyStock `json:"stocks"`
	PiecesBySemaforizationColor map[string]int  `json:"pieces_by_semaforization_color"`
	TotalPieces                 int16           `json:"total_pieces"`
}

func (pharmacyStocksDetailed *PharmacyStocksDetailed) CountAndCategorizePieces() {
	var totalPieces int
	multipleSemaforizationCounter := make(map[string]int, 0)
	for _, stock := range pharmacyStocksDetailed.Stocks {
		multipleSemaforizationCounter[string(stock.SemaforizationColor)] += int(stock.Pieces)
	}
	for _, pieces := range multipleSemaforizationCounter {
		totalPieces += pieces
	}
	pharmacyStocksDetailed.PiecesBySemaforizationColor = multipleSemaforizationCounter
	pharmacyStocksDetailed.TotalPieces = int16(totalPieces)
}
