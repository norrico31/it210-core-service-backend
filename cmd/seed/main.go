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
	defer db.Close()
	seeders.SeedRoles(db)
	seeders.SeedPriorities(db)
	seeders.SeedSegments(db)
	seeders.SeedUsers(db)
	seeders.SeedStatuses(db)
	seeders.SeedProjects(db)
	seeders.SeedWorkspace(db)
	seeders.SeedTasks(db)
	seeders.SeedProjectTasks(db)
	log.Println("Seeding successfully complete.")
}
