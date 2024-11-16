package main

import (
	"log"

	"github.com/norrico31/it210-core-service-backend/cmd/seed/seeders"
	"github.com/norrico31/it210-core-service-backend/db"
)

func main() {
	db, err := db.NewPostgresStorage()
	if err != nil {
		log.Fatalf("Failed to connect to the database %v", err)
	}
	// db.Exec(`DROP DATABASE it210`)
	// db.Exec(`CREATE DATABASE it210`)
	seeders.SeedRoles(db)
	seeders.SeedUsers(db)
	seeders.SeedStatuses(db)
	seeders.SeedProjects(db)
	seeders.SeedTasks(db)
	log.Println("Seeding successfully complete.")
}
