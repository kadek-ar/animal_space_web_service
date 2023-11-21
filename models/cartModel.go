package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	Note     string
	Image    string
	Quantity int
	UserID   int
	User     User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AnimalID int
	Animal   Animal `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
