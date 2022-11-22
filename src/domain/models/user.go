package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        string     `gorm:"primaryKey" json:"id,omitempty"`
	RoleId    string     `json:"role_id"`
	Role      *Role      `gorm:"foreignKey:role_id" json:"role,omitempty"`
	Name      string     `json:"name"`
	Surname   string     `json:"surname"`
	FullName  string     `json:"full_name,omitempty"`
	Email     string     `json:"email,omitempty" validate:"required,email"`
	Password  string     `json:"password,omitempty" validate:"required,min=8"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time  `gorm:"index" json:"deleted_at"`
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
	var numOfRecord int64
	dateUnixStr := fmt.Sprint(time.Now().Unix())
	dateUnixSufix := dateUnixStr[len(dateUnixStr)-3:]
	tx.Model(User{}).Count(&numOfRecord)
	numOfRecord += 100
	user.Id = fmt.Sprintf("CAN%v%v", numOfRecord, dateUnixSufix)
	user.Password = hashPassword(user.Password)
	return
}

func (user *User) RestPassword(password string) {
	user.Password = hashPassword(password)
}

func (user *User) WithoutPassword() User {
	user.Password = ""
	return *user
}

func (user *User) AfterFind(tx *gorm.DB) error {
	user.FullName = user.Name + " " + user.Surname
	return nil
}
