package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type customLogger struct{}

func (l *customLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *customLogger) Verbose() bool {
	return true
}

func main() {
	fmt.Println("--Migrate Start--")

	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	dbURL := os.Getenv("POSTGRES_URI")
	m, err := migrate.New(
		"file://../migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	m.Log = &customLogger{}

	// UP
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations ran successfully")

	// ROLLBACK
	// err = m.Down()
	// if err != nil && err != migrate.ErrNoChange {
	// 	log.Fatalf("Failed to roll back migrations: %v", err)
	// }
	// log.Println("Migrations rolled back successfully")
}
