package datasources

import (
	"cecan_inventory/domain/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc DepartmentDataSource) CreateDepartment(department models.Department) (uuid.UUID, error) {
	err := dataSrc.DbPsql.Create(&department).Error
	if err != nil {
		return uuid.Nil, err
	}
	return department.Id, nil
}

func (dataSrc DepartmentDataSource) GetDepartments(includeDeleted bool, limit int, offset int) ([]models.Department, error) {
	var departments = make([]models.Department, 0)
	err := dataSrc.DbPsql.Raw(
		fmt.Sprintf("SELECT * FROM get_departments_ordered_by_floor(%v,%v,%v);", includeDeleted, limit, offset),
	).Scan(&departments).Error
	if err != nil {
		return departments, err
	}
	return departments, nil
}

func (dataSrc DepartmentDataSource) GetDepartmentById(id string) (models.DepartmentDetailed, error) {
	var department models.DepartmentDetailed
	err := dataSrc.DbPsql.
		Table("departments").
		Unscoped().
		Where("id = ?", id).
		Preload("ResponsibleUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "surname", "role_id")
		}).
		Find(&department).
		Error
	if err != nil {
		return department, err
	}
	return department, nil
}

func (dataSrc DepartmentDataSource) GetDepartmentByName(name string) (models.DepartmentDetailed, error) {
	var department models.DepartmentDetailed
	err := dataSrc.DbPsql.
		Table("departments").
		Unscoped().
		Where("name = ?", name).
		Preload("ResponsibleUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "surname", "role_id")
		}).
		Find(&department).
		Error
	if err != nil {
		return department, err
	}
	return department, nil
}

func (dataSrc DepartmentDataSource) UpdateDepartment(id string, department models.Department) error {
	err := dataSrc.DbPsql.
		Omit("id", "created_at", "responsible_user_id").
		Where("id = ?", id).
		Updates(&department).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc DepartmentDataSource) AssingResponsibleToDepartment(id string, userId string) error {
	err := dataSrc.DbPsql.
		Model(&models.Department{}).
		Where("id = ?", id).
		Update("responsible_user_id", userId).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc DepartmentDataSource) ReactivateDepartment(id string) error {
	err := dataSrc.DbPsql.
		Table("departments").
		Where("id = ?", id).
		Update("deleted_at", nil).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc DepartmentDataSource) DeleteDepartment(id string) error {
	errInTransaction := dataSrc.DbPsql.Transaction(func(tx *gorm.DB) error {
		errUpdating := tx.
			Model(&models.Department{}).
			Where("id = ?", id).
			Update("responsible_user_id", nil).
			Error
		if errUpdating != nil {
			return errUpdating
		}
		errDeleting := tx.
			Where("id = ?", id).
			Delete(&models.Department{}).
			Error
		if errDeleting != nil {
			return errDeleting
		}
		return nil
	})
	if errInTransaction != nil {
		return errInTransaction
	}
	return nil
}
