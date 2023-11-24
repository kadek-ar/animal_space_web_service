package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Note         string
	Status       string
	Total        int
	NumberOfItem int
	UserID       int
	User         User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Animal       []Animal `gorm:"many2many:transaction_animals;"`
}

type TransactionAnimal struct {
	AnimalID      int  `gorm:"primaryKey"`
	TransactionID uint `gorm:"primaryKey"`
	Note          string
	Quantity      int
	Price         int
	Images        string
}

type GetTransactionDetail struct {
	TransactionID     int
	AnimalID          int
	Images            string
	Quantity          int
	AnimalName        string
	AnimalGender      string
	AnimalType        string
	AnimalDescription string
	AnimalImage       string
	AnimalPrice       int
	AnimalCategory    string
	ShelterID         int
	ShelterName       string
	ShelterPhone      string
}

type GetShelterTransaction struct {
	TransactionID int       `json:"transaction_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	AnimalCount   int       `json:"animal_count"`
	Total         int       `json:"total_price"`
	ShelterID     int       `json:"shelter_id"`
}
