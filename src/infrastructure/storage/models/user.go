package models

import (
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	RoleId    string    `json:"role_id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at"`
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println("User's password couldn't be hashed!")
	}
	return string(bytes)
}

func (user User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Password = hashPassword(user.Password)
	return
}

func (user User) ToJSON() iris.Map {
	return iris.Map{
		"id":      user.Id,
		"role_id": user.RoleId,
		"name":    user.Name,
		"surname": user.Surname,
		"email":   user.Email,
	}

}
