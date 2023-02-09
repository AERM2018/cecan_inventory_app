package datasources

import (
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	"errors"

	"gorm.io/gorm"
)

type StorehouseUtilitiesDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc StorehouseUtilitiesDataSource) CreateStorehouseUtility(utility models.StorehouseUtility) error {
	err := dataSrc.DbPsql.Create(&utility).Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc StorehouseUtilitiesDataSource) GetStorehouseUtilities(filters models.StorehouseUtilitiesFilters) ([]models.StorehouseUtilityDetailed, int, error) {
	var totalRecords int64
	totalPages := 1
	conditionStringFromJson := common.StructJsonSerializer(models.StorehouseUtilitiesFilters{
		UtilityKey:  filters.UtilityKey,
		UtilityName: filters.UtilityName,
	}, "json", "OR")
	storehouseUtilities := make([]models.StorehouseUtilityDetailed, 0)
	dbPointer := dataSrc.DbPsql.Model(&models.StorehouseUtility{})
	if filters.IncludeDeleted {
		dbPointer = dbPointer.Unscoped()
	}
	dbPointer = dbPointer.
		Preload("Presentation", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "updated_at", "deleted_at")
		}).
		Preload("Unit", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "updated_at", "deleted_at")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "updated_at", "deleted_at")
		}).
		Where(conditionStringFromJson).
		Count(&totalRecords).
		Limit(filters.Limit).
		Offset((filters.Page - 1) * filters.Limit).
		Find(&storehouseUtilities)

	totalPages = int(totalRecords) / filters.Limit
	if totalPages*filters.Limit != int(totalRecords) {
		totalPages += 1
	}
	if dbPointer.Error != nil {
		return storehouseUtilities, 0, dbPointer.Error
	}
	return storehouseUtilities, totalPages, nil
}

func (dataSrc StorehouseUtilitiesDataSource) GetStorehouseUtilityByKey(key string) (models.StorehouseUtilityDetailed, error) {
	var utilityDetailed models.StorehouseUtilityDetailed
	err := dataSrc.DbPsql.Unscoped().Model(&models.StorehouseUtility{}).
		Preload("Presentation", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "updated_at", "deleted_at")
		}).
		Preload("Unit", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "updated_at", "deleted_at")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "updated_at", "deleted_at")
		}).
		Where("key = ?", key).
		Take(&utilityDetailed).Error

	if err != nil {
		return utilityDetailed, err
	}

	return utilityDetailed, nil
}

func (dataSrc StorehouseUtilitiesDataSource) UpdateStorehouseUtility(key string, utility models.StorehouseUtility) (string, error) {
	err := dataSrc.DbPsql.
		Model(&models.StorehouseUtility{}).
		Where("key = ?", key).
		Updates(&utility).Error
	if err != nil {
		return "", err
	}
	return utility.Key, nil
}

func (dataSrc StorehouseUtilitiesDataSource) ReactivateStorehouseUtility(key string) error {
	err := dataSrc.DbPsql.
		Model(&models.StorehouseUtility{}).
		Unscoped().
		Where("key = ?", key).
		Update("deleted_at", nil).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc StorehouseUtilitiesDataSource) DeleteStorehouseUtility(key string) error {
	err := dataSrc.DbPsql.
		Where("key = ?", key).
		Delete(&models.StorehouseUtility{}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc StorehouseUtilitiesDataSource) GetStorehouseUtilityCategories() ([]models.StorehouseUtilityCategory, error) {
	storehouseUtilityCategories := make([]models.StorehouseUtilityCategory, 0)
	err := dataSrc.DbPsql.Find(&storehouseUtilityCategories).Error
	if err != nil {
		return storehouseUtilityCategories, err
	}
	return storehouseUtilityCategories, nil
}

func (dataSrc StorehouseUtilitiesDataSource) GetStorehouseUtilityPresentations() ([]models.StorehouseUtilityPresentation, error) {
	storehouseUtilityPresentations := make([]models.StorehouseUtilityPresentation, 0)
	err := dataSrc.DbPsql.Find(&storehouseUtilityPresentations).Error
	if err != nil {
		return storehouseUtilityPresentations, err
	}
	return storehouseUtilityPresentations, nil
}

func (dataSrc StorehouseUtilitiesDataSource) GetStorehouseUtilityUnits() ([]models.StorehouseUtilityUnit, error) {
	storehouseUtilityUnits := make([]models.StorehouseUtilityUnit, 0)
	err := dataSrc.DbPsql.Find(&storehouseUtilityUnits).Error
	if err != nil {
		return storehouseUtilityUnits, err
	}
	return storehouseUtilityUnits, nil
}

func (dataSrc StorehouseUtilitiesDataSource) IsStockFromUtilityUsed(key string) bool {
	var utilityStock models.StorehouseStock
	err := dataSrc.DbPsql.
		Where("storehouse_utility_key = ? and quantity_parsed_used > ?", key, 0).
		Take(&utilityStock).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
