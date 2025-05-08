package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
		log.Printf("Warning: Error loading .env file: %v", err)
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

	currentVersion, dirty, verErr := m.Version()
	if verErr == nil {
		log.Printf("Current migration version before operation: %d (dirty: %t)", currentVersion, dirty)
	} else if verErr != migrate.ErrNilVersion {
		log.Printf("Error getting current version: %v", verErr)
	} else {
		log.Printf("No migrations have been applied yet")
	}

	args := os.Args
	if len(args) < 3 {
		log.Fatalf("Usage: go run main.go [up|down] [target_version]")
	}

	command := args[1]
	targetVersionStr := args[2]

	switch command {
	case "up":
		if targetVersionStr == "latest" {
			log.Println("Migrating up to the latest version...")
			err = m.Up()
			if err != nil {
				if err == migrate.ErrNoChange {
					log.Printf("No change: Database already at the latest version")
				} else {
					log.Fatalf("Failed to migrate up: %v", err)
				}
			} else {
				log.Println("Successfully migrated to the latest version")
			}
		} else {
			targetVersion, err := strconv.Atoi(targetVersionStr)
			if err != nil {
				log.Fatalf("Invalid target version: %v", err)
			}

			log.Printf("Migrating up to version %d...", targetVersion)
			err = m.Migrate(uint(targetVersion))
			if err != nil {
				if err == migrate.ErrNoChange {
					log.Printf("No change: Database already at version %d", targetVersion)
				} else {
					log.Fatalf("Failed to migrate up to version %d: %v", targetVersion, err)
				}
			} else {
				log.Printf("Successfully migrated up to version %d", targetVersion)
			}
		}

	case "down":
		if targetVersionStr == "0" {
			log.Println("Rolling back all migrations...")
			err = m.Down()
			if err != nil {
				if err == migrate.ErrNoChange {
					log.Printf("No change: No migrations to roll back")
				} else {
					log.Fatalf("Failed to roll back all migrations: %v", err)
				}
			} else {
				log.Println("Successfully rolled back all migrations")
			}
		} else if targetVersionStr == "prev" {
			log.Println("Rolling back one migration...")
			err = m.Steps(-1)
			if err != nil {
				if err == migrate.ErrNoChange {
					log.Printf("No change: No migrations to roll back")
				} else {
					log.Fatalf("Failed to roll back one migration: %v", err)
				}
			} else {
				log.Println("Successfully rolled back one migration")
			}
		} else {
			targetVersion, err := strconv.Atoi(targetVersionStr)
			if err != nil {
				log.Fatalf("Invalid target version: %v", err)
			}

			currentVersion, _, _ := m.Version()
			if uint(targetVersion) >= currentVersion {
				log.Fatalf("Target version (%d) must be less than current version (%d) for down migration",
					targetVersion, currentVersion)
			}

			log.Printf("Rolling back to version %d...", targetVersion)
			err = m.Migrate(uint(targetVersion))
			if err != nil {
				if err == migrate.ErrNoChange {
					log.Printf("No change: Database already at version %d", targetVersion)
				} else {
					log.Fatalf("Failed to roll back to version %d: %v", targetVersion, err)
				}
			} else {
				log.Printf("Successfully rolled back to version %d", targetVersion)
			}
		}

	default:
		log.Fatalf("Invalid command: %s. Use 'up' or 'down'", command)
	}

	finalVersion, dirty, verErr := m.Version()
	if verErr == nil {
		log.Printf("Current migration version after operation: %d (dirty: %t)", finalVersion, dirty)
	} else if verErr != migrate.ErrNilVersion {
		log.Printf("Error getting current version: %v", verErr)
	} else {
		log.Printf("No migrations are applied")
	}

	fmt.Println("--Migrate Complete--")
}
