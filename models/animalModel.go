package models

import (
	"time"

	"gorm.io/gorm"
)

type Animal struct {
	gorm.Model
	Name         string
	Image        string
	Gender       string
	Type         string
	Age          int
	Description  string
	Quantity     int
	Status       string
	Price        int
	CategoriesID int
	Categories   Categories `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ShelterID    int
	Shelter      Shelter `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type GetAllAnimal struct {
	Id            int
	Name          string
	Type          string
	Age           int
	Description   string
	Quantity      int
	Price         int
	Status        string
	Image         string
	Category_id   int
	Category_name string
	Shelter_id    int
	Shelter_name  string
	Created_at    time.Time
}
