package migrators

import (
	"blogpost/models"

	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) {
	db.AutoMigrate(&models.LookUp{})
}
