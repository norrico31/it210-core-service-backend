package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	_ "github.com/golang-migrate/migrate/v4/source/file" // File source driver
	_ "github.com/lib/pq"                                // PostgreSQL driver

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/norrico31/it210-core-service-backend/db"
)

func main() {
	db, err := db.NewPostgresStorage()

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations up: %v", err)
		} else {
			log.Println("Migrations applied successfully.")
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations down: %v", err)
		} else {
			log.Println("Drop tables applied successfully.")
		}
	}
}
