package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var stringConnection = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", "localhost", "root", "postgresql", "cecan_dev", "5432")
var (
	DBInstance *gorm.DB
	err        error
)

func GetDbInstance() (*gorm.DB, error) {
	if DBInstance != nil {
		DBInstance, err = gorm.Open(postgres.Open(stringConnection), &gorm.Config{})
	}
	return DBInstance, err
}
