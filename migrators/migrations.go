package migrators

import (
	"blogpost/models"

	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comments{}, models.Views{})
}
