package seeders

import (
	"database/sql"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedWorkspace(db *sql.DB) error {
	// get projectId
	projects := make(map[string]int)
	projectRows, err := db.Query(`
		SELECT 
			id, name
		FROM projects
	`)

	if err != nil {
		log.Printf("Failed to fetch projects: %v\n", err)
		return err
	}
	defer projectRows.Close()

	for projectRows.Next() {
		var id int
		var name string
		if err = projectRows.Scan(&id, &name); err != nil {
			log.Printf("Failed to scan projects: %v\n", err)
			return err
		}
		projects[name] = id
	}

	interpreter := projects["Project 124 Interpreter"]
	dostv := projects["Project 210 Web App DOSTV"]

	workspaces := []entities.Workspace{
		{
			Name:        "PENDING",
			Description: "low level",
			ProjectID:   interpreter,
		},
		{
			Name:        "PENDING",
			Description: "high level",
			ProjectID:   dostv,
		},
		{
			Name:        "ONGOING",
			Description: "medium level",
			ProjectID:   interpreter,
		},
		{
			Name:        "ONGOING",
			Description: "critical level",
			ProjectID:   dostv,
		},
	}

	for _, workspace := range workspaces {
		_, err := db.Exec(`
			INSERT INTO workspaces (name, description, projectId, createdAt, updatedAt)
			VALUES ($1, $2, $3, $4, $5)
		`, workspace.Name, workspace.Description, workspace.ProjectID, time.Now(), time.Now())

		if err != nil {
			log.Printf("Failed to insert workspace %s: %v\n", workspace.Name, err)
			return err
		}
		log.Printf("Successfully inserted workspace %s\n", workspace.Name)
	}

	return nil
}
