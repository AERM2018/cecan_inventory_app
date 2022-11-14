package models

import "github.com/google/uuid"

type FixedAssetsItemsRequests struct {
	Id                   uuid.UUID  `gorm:"primaryKey;default:'uuid_generate_v4();'" json:"id"`
	FixedAssetKey        string     `gorm:"foreignKey:fixed_asset_key" json:"fixed_asset_key"`
	FixedAsset           FixedAsset `json:"details"`
	FixedAssetsRequestId uuid.UUID  `gorm:"foreignKey:fixed_assets_requet_id" json:"fixed_assets_requet_id"`
}
