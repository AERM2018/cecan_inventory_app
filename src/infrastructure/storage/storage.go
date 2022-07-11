package storage

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"cecan_inventory/src/infrastructure/storage/migrator"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var stringConnection = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", "192.168.0.8", "postgres", "C4s10*2016", "cecan_dev", "5432")
var (
	DBInstance *gorm.DB = nil
	err        error
)

func Connect() (*gorm.DB, error) {
	if DBInstance == nil {
		DBInstance, err = gorm.Open(postgres.Open(stringConnection), &gorm.Config{})
	}
	if err := Migrate(DBInstance, true); err != nil {
		return DBInstance, err
	}
	return DBInstance, err
}

func Migrate(db *gorm.DB, isMigrationUp bool) error {
	psql, err := db.DB();
	if err != nil{
		return nil;
	}
	driver, err := migrator.GetPsqlDriver(psql);
	if err != nil{
		return err
	}
	return migrator.Exec(driver,os.Getenv("PSQL_DB_NAME"),getMigrationsPath(), isMigrationUp);
}

func getMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "migrations")
}