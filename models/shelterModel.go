package models

import (
	"gorm.io/gorm"
)

type Shelter struct {
	gorm.Model
	Name        string
	Phone       string
	Description string
	Address     string
	Status      string
	Note        string
	UserID      int
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type GetShelter struct {
	Id          int
	Name        string
	Phone       string
	Description string
	Address     string
	Status      string
	User_id     int
	Email_user  string
	Owner_name  string
}

type Categories struct {
	gorm.Model
	Name  string
	Image string
}
