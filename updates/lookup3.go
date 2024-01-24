package migrators

import (
	"blogpost/models"
	"fmt"
)

func (u *LookUpDb) Lookup3() {
	fmt.Println("lookup3 called")
	u.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comments{}, models.Views{}, models.Category{})
}
