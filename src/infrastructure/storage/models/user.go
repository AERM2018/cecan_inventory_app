package models

import "time"

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
