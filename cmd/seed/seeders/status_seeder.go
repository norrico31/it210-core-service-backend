package seeders

import (
	"database/sql"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedStatuses(db *sql.DB) error {
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
	}

	for _, status := range statuses {
		_, err := db.Exec(`
			INSERT INTO statuses (name, description, createdAt, updatedAt)
			VALUES ($1, $2, $3, $4)
		`, status.Name, status.Description, time.Now(), time.Now())

		if err != nil {
			log.Printf("Failed to insert status %s: %v\n", status.Name, err)
			return err
		}
		log.Printf("Successfully inserted status %s\n", status.Name)
	}

	return nil
}
