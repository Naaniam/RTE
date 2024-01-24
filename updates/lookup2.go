package migrators

import (
	"blogpost/models"
)

func (u *LookUpDb) Lookup2() {
	u.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comments{}, models.Views{})
}
