// cmd/seed/seeders/project_seeder.go
package seeders

import (
	"database/sql"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedProjects(db *sql.DB) error {
	// Sample data for seeding
	projects := []entities.Project{
		{
			Name:        "Project Alpha",
			Description: "Description for Project Alpha",
			Users:       []entities.User{{ID: 1}, {ID: 2}}, // Assumed existing user IDs
			Tasks:       []entities.Task{{ID: 1}, {ID: 2}}, // Assumed existing task IDs
		},
		{
			Name:        "Project Beta",
			Description: "Description for Project Beta",
			Users:       []entities.User{{ID: 3}, {ID: 4}}, // Assumed existing user IDs
			Tasks:       []entities.Task{{ID: 3}, {ID: 4}}, // Assumed existing task IDs
		},
		// Add more projects as needed
	}

	// Insert each project into the database
	for _, project := range projects {
		// Handle nullable lastActiveAt (set to nil for now)
		var lastActiveAt *time.Time

		// Insert project into the database
		_, err := db.Exec(`
			INSERT INTO projects (name, description, lastActiveAt, createdAt, updatedAt)
			VALUES ($1, $2, $3, $4, $5)
		`, project.Name, project.Description, lastActiveAt, time.Now(), time.Now())

		if err != nil {
			log.Printf("Failed to insert project %s: %v\n", project.Name, err)
			return err
		}
		log.Printf("Successfully inserted project %s\n", project.Name)

		// Fetch the inserted project ID to associate users and tasks
		var projectId int
		err = db.QueryRow(`SELECT id FROM projects WHERE name = $1`, project.Name).Scan(&projectId)
		if err != nil {
			log.Printf("Failed to fetch project ID for %s: %v\n", project.Name, err)
			return err
		}

		// Insert project-user associations
		for _, user := range project.Users {
			_, err := db.Exec(`
				INSERT INTO users_projects (project_id, user_id) 
				VALUES ($1, $2)
			`, projectId, user.ID)

			if err != nil {
				log.Printf("Failed to associate user %d with project %d: %v\n", user.ID, projectId, err)
				return err
			}
			log.Printf("Successfully associated user %d with project %d\n", user.ID, projectId)
		}

		// Insert project-task associations
		for _, task := range project.Tasks {
			_, err := db.Exec(`
				INSERT INTO project_tasks (project_id, task_id)
				VALUES ($1, $2)
			`, projectId, task.ID)

			if err != nil {
				log.Printf("Failed to associate task %d with project %d: %v\n", task.ID, projectId, err)
				return err
			}
			log.Printf("Successfully associated task %d with project %d\n", task.ID, projectId)
		}
	}

	return nil
}
