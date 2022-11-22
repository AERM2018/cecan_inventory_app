package datasources

import (
	"cecan_inventory/domain/models"

	"gorm.io/gorm"
)

type PasswordResetCodesDataSource struct {
	DbPsql *gorm.DB
}

func (dataSrc PasswordResetCodesDataSource) CreateCode(credentialsRestart models.PasswordResetCode) error {
	err := dataSrc.DbPsql.Create(&credentialsRestart).Error
	if err != nil {
		return err
	}
	return nil
}

func (dataSrc PasswordResetCodesDataSource) GetCode(code string) (models.PasswordResetCode, error) {
	var codeFound models.PasswordResetCode
	err := dataSrc.DbPsql.Where("code = ?", code).Take(&codeFound).Error
	if err != nil {
		return codeFound, err
	}
	return codeFound, nil
}

func (dataSrc PasswordResetCodesDataSource) UseCode(code string) error {
	err := dataSrc.DbPsql.Table("password_reset_codes").Where("code = ?").UpdateColumn("is_used", true).Error
	if err != nil {
		return err
	}
	return nil
}
