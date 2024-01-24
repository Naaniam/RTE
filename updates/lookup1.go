package migrators

import (
	"blogpost/models"

	"gorm.io/gorm"
)

type LookUpDb struct {
	DB *gorm.DB
}

func NewLookUpDB(db *gorm.DB) *LookUpDb {
	return &LookUpDb{DB: db}
}

func (u *LookUpDb) Lookup1() {
	u.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comments{}, models.Views{})
}
