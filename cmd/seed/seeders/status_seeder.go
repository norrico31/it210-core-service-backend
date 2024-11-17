package seeders

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedStatuses(db *sql.DB) error {
	// Sample data for seeding
	statuses := []entities.Status{
		{
			Name:        "Active",
			Description: "The status is active and in progress.",
		},
		{
			Name:        "Completed",
			Description: "The task or project is completed.",
		},
		{
			Name:        "Pending",
			Description: "The task or project is pending and waiting to be started.",
		},
		{
			Name:        "Archived",
			Description: "The task or project is archived and no longer active.",
		},
		{
			Name:        "Not Started",
			Description: "The task or project is not yet started.",
		},
		{
			Name:        "In Progress",
			Description: "The task or project is currently in progress.",
		},
		// Add more statuses as needed
	}

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	for _, status := range statuses {
		// Increment the counter for the WaitGroup
		wg.Add(1)

		// Use a goroutine to insert each status concurrently
		go func(status entities.Status) {
			defer wg.Done() // Decrement the counter when the goroutine completes

			// Insert status into the database
			_, err := db.Exec(`
				INSERT INTO statuses (name, description, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4)
			`, status.Name, status.Description, time.Now(), time.Now())

			if err != nil {
				log.Printf("Failed to insert status %s: %v\n", status.Name, err)
				return
			}
			log.Printf("Successfully inserted status %s\n", status.Name)
		}(status)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}
