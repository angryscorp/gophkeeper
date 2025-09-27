package main

import (
	"gophkeeper/client/internal/config"
	"log"
	"os"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

const cfgFileName = "config.json"

func main() {
	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// Load config from a file
	cfg, err := config.LoadFromFile(execPath, cfgFileName)
	if err != nil {
		log.Fatal(err)
	}

	// Bootstrap
	program, closes := bootstrap(cfg)
	defer func() {
		for _, f := range closes {
			f()
		}
	}()

	// Run the program
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
