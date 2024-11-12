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
			// Users:       []entities.User{{ID: 1}, {ID: 2}}, // Assumed existing user IDs
			// Tasks:       []entities.Task{{ID: 1}, {ID: 2}}, // Assumed existing task IDs
		},
		{
			Name:        "Project Beta",
			Description: "Description for Project Beta",
			// Users:       []entities.User{{ID: 3}, {ID: 4}}, // Assumed existing user IDs
			// Tasks:       []entities.Task{{ID: 3}, {ID: 4}}, // Assumed existing task IDs
		},
		// Add more projects as needed
	}

	// Insert each project into the database
	for _, project := range projects {
		// Insert project into the database
		_, err := db.Exec(`
			INSERT INTO projects (name, description, createdAt, updatedAt)
			VALUES ($1, $2, $3, $4)
		`, project.Name, project.Description, time.Now(), time.Now())

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

		// Ensure all users exist before associating them with the project
		for _, user := range project.Users {
			// Check if the user exists in the users table
			var userId int
			err := db.QueryRow(`SELECT id FROM users WHERE id = $1`, user.ID).Scan(&userId)
			if err != nil {
				log.Printf("User ID %d does not exist in users table: %v\n", user.ID, err)
				continue // Skip this user if not found
			}

			// Insert project-user associations
			_, err = db.Exec(`
				INSERT INTO users_projects (project_id, user_id) 
				VALUES ($1, $2)
			`, projectId, userId)

			if err != nil {
				log.Printf("Failed to associate user %d with project %d: %v\n", userId, projectId, err)
				return err
			}
			log.Printf("Successfully associated user %d with project %d\n", userId, projectId)
		}

		// Insert project-task associations
		for _, task := range project.Tasks {
			// Check if the task exists in the tasks table
			var taskId int
			err := db.QueryRow(`SELECT id FROM tasks WHERE id = $1`, task.ID).Scan(&taskId)
			if err != nil {
				log.Printf("Task ID %d does not exist in tasks table: %v\n", task.ID, err)
				continue // Skip this task if not found
			}

			_, err = db.Exec(`
				INSERT INTO project_tasks (project_id, task_id)
				VALUES ($1, $2)
			`, projectId, taskId)

			if err != nil {
				log.Printf("Failed to associate task %d with project %d: %v\n", taskId, projectId, err)
				return err
			}
			log.Printf("Successfully associated task %d with project %d\n", taskId, projectId)
		}
	}

	return nil
}
