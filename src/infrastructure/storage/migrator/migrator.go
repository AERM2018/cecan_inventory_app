package migrator

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)
func Exec(driver database.Driver, dbName string, migrationsPath string, isMigrationsUp bool) error {
	iMigrate, err := migrate.NewWithDatabaseInstance(
		"file:"+migrationsPath,
		dbName,
		driver,
	)
	if err != nil{
		return err;
	}
	if isMigrationsUp {
		if err := iMigrate.Up(); err != nil {
			return err
		}
	} else {
		if err := iMigrate.Down(); err != nil {
			return err
		}
	}
	return nil;
}

func GetPsqlDriver(db *sql.DB) (database.Driver, error) {
	return postgres.WithInstance(db, &postgres.Config{});
}