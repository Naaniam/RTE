package main

import (
	driver "blogpost/drivers"
	"blogpost/migrators"
	"blogpost/repository"
	"blogpost/router"
)

func main() {
	dbConnection := driver.SQLDriver()
	migrators.Migrations(dbConnection)
	router.Routing(repository.NewDbConnection(dbConnection))
}
