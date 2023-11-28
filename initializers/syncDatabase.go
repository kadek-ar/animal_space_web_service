package initializers

import "project/web-service-gin/models"

func SyncDatabase() {
	DB.AutoMigrate(
		&models.User{},
		&models.Shelter{},
		&models.Category{},
		&models.Animal{},
		&models.Cart{},
		&models.Transaction{},
		&models.TransactionAnimal{},
		&models.Banner{},
	)
}
