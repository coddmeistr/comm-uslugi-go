package initializers

import "golang-uslugi-server/m/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Request{})
	DB.AutoMigrate(&models.Worker{})
}
