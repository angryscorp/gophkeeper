package main

import (
	"fmt"

	"gophkeeper/server/internal/config"
	"gophkeeper/server/internal/repository/migration"
	usersrepo "gophkeeper/server/internal/repository/users/impl"
)

func main() {
	fmt.Println("hello world - server")
	
	cfg := config.Config{
		DatabaseDSN: "",
		Debug:       true,
	}

	if err := migration.MigratePostgres(cfg.DatabaseDSN); err != nil {
		panic(err)
	}

	repository, err := usersrepo.New(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	fmt.Println(repository)
}
