package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Username string
	Role     string
	Password string `gorm:"size:100"`
	Status   string
	Hash     string
}

type GetUser struct {
	Email    string `gorm:"unique"`
	Username string
	Role     string
}
