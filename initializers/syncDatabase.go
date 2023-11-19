package initializers

import "project/web-service-gin/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Shelter{}, &models.Categories{}, &models.Animal{})
}
