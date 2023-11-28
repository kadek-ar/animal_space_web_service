package models

import "gorm.io/gorm"

type Banner struct {
	gorm.Model
	Image string `json:"image"`
}
