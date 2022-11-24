package models

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	AccessCredentials struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password"`
	}

	AccessCredentialsRestart struct {
		Email       string `json:"email,omitempty"`
		OldPassword string `json:"old_password,omitempty"`
		NewPassword string `json:"new_password,omitempty"`
		ResetCode   string `json:"code,omitempty"`
	}

	PasswordResetCode struct {
		Id        uuid.UUID `gorm:"primaryKey;default:'uuid_generate_v4();'"`
		Code      string
		IsUsed    bool
		CreatedAt time.Time
		ExpiresAt time.Time
	}
)

func (passwordCode PasswordResetCode) HashToken() string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("PASSWORD_SECRET")), 10)
	if err != nil {
		fmt.Println("Password reset code couldn't be hashed!")
	}
	return string(bytes)
}

func (passwordCode *PasswordResetCode) CheckToken(token string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(token), []byte(os.Getenv("PASSWORD_SECRET")))
	return err == nil
}
