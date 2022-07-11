package datasources

import "gorm.io/gorm"

type UserDataSource struct {
	psql *gorm.DB
}

func (dataSrc *UserDataSource) login() {
	// dataSrc.psql.Model()
}