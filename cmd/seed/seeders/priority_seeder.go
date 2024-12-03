package seeders

import (
	"database/sql"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedPriorities(db *sql.DB) error {
	priorities := []entities.Priority{
		{
			Name:        "Low",
			Description: "low level",
		},
		{
			Name:        "Medium",
			Description: "medium level",
		},
		{
			Name:        "High",
			Description: "high level",
		},
		{
			Name:        "Critical",
			Description: "critical level",
		},
	}

	for _, priority := range priorities {
		_, err := db.Exec(`
				INSERT INTO priorities (name, description, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4)
			`, priority.Name, priority.Description, time.Now(), time.Now())

		if err != nil {
			log.Printf("Failed to insert priority %s: %v\n", priority.Name, err)
			return err
		}
		log.Printf("Successfully inserted priority %s\n", priority.Name)
	}

	return nil
}
