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