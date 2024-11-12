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

	seeders.SeedStatuses(db)
	seeders.SeedRoles(db)
	seeders.SeedTasks(db)
	seeders.SeedProjects(db)
	log.Println("Seeding successfully complete.")
}
