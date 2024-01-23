package driver

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SQLDriver() *gorm.DB {
	var err error

	dsn := "root:password@tcp(localhost:3306)/blogpost?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection Established")
	return db
}
