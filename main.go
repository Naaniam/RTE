package main

import (
	driver "blogpost/drivers"
	"blogpost/migrators"
	"blogpost/repository"
	"blogpost/router"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// Open or create a log file
	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		//	logger.Println("Error opening log file:", err)
		fmt.Println("Error opening log file:", err)
		return
	}

	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)
	logger := log.New(multiWriter, "Blog-Post ", log.LstdFlags)

	dbConnection := driver.SQLDriver()
	migrators.Migrations(dbConnection)
	router.Routing(repository.NewDbConnection(dbConnection, logger))
}
