package models

import "github.com/google/uuid"

type FixedAssetsItemsRequests struct {
	Id                   uuid.UUID  `gorm:"primaryKey;default:'uuid_generate_v4();'" json:"id"`
	FixedAssetId         uuid.UUID  `gorm:"foreignKey:fixed_asset_id" json:"fixed_asset_id"`
	FixedAsset           FixedAsset `json:"details"`
	FixedAssetsRequestId uuid.UUID  `gorm:"foreignKey:fixed_assets_request_id" json:"fixed_assets_request_id"`
}
