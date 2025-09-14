package main

import (
	"fmt"
	"time"

	"gophkeeper/server/internal/config"
	"gophkeeper/server/internal/repository/migration"
	usersrepo "gophkeeper/server/internal/repository/users/impl"
)

func main() {
	fmt.Println("hello world - server")

	cfg := config.Config{
		DatabaseDSN: "postgres://db_user:db_password@postgres:5432/gophkeeper?sslmode=disable",
		Debug:       true,
	}

	if err := migration.MigratePostgres(cfg.DatabaseDSN); err != nil {
		panic(err)
	}

	repository, closeDB, err := usersrepo.New(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	defer closeDB()

	fmt.Println(repository)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("tick - %s\n", time.Now())
	}

}
