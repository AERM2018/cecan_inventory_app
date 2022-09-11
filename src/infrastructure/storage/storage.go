package storage

import (
	"cecan_inventory/infrastructure/storage/migrator"
	"cecan_inventory/infrastructure/storage/seeds"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBInstance *gorm.DB = nil
	err        error
)

func Connect() (*gorm.DB, error) {
	var stringConnection = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("CECAN_DB_HOST"),
		os.Getenv("CECAN_DB_USER"),
		os.Getenv("CECAN_DB_PASSWD"),
		os.Getenv("CECAN_DB_NAME"),
		"5432") // Get stringConnection with help of the env file
	if DBInstance == nil {
		// 		newLogger := logger.New(
		//   log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		//   logger.Config{
		//     SlowThreshold:              time.Second,   // Slow SQL threshold
		//     LogLevel:                   logger.Silent, // Log level
		//     IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
		//     Colorful:                  false,          // Disable color
		//   },
		// )
		// DBInstance, err = gorm.Open(postgres.Open(stringConnection), &gorm.Config{Logger: newLogger})
		DBInstance, err = gorm.Open(postgres.Open(stringConnection))
	}
	err := Migrate(DBInstance, true)
	if os.Getenv("SEED_APPLIED") == "false" {
		for _, seed := range seeds.All() {
			if err := seed.Run(DBInstance); err != nil {
				log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
			}
		}
		os.Setenv("SEED_APPLIED", "true")
	}
	if err != nil {
		return DBInstance, err
	}
	return DBInstance, nil
}

func PruneData(DBInstance *gorm.DB) (*gorm.DB, error) {
	fmt.Println("Prunning data...")
	err := Migrate(DBInstance, false)
	if err != nil {
		return DBInstance, err
	}
	return DBInstance, nil
}

func Migrate(db *gorm.DB, isMigrationUp bool) error {
	psql, err := db.DB()
	if err != nil {
		return nil
	}
	driver, err := migrator.GetPsqlDriver(psql)
	if err != nil {
		return err
	}
	return migrator.Exec(driver, os.Getenv("PSQL_DB_NAME"), getMigrationsPath(), isMigrationUp)
}

func getMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "migrations")
}
