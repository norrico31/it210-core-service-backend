package seeders

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedProjects(db *sql.DB) error {
	projects := []entities.Project{
		{
			Name:        "Project Alpha",
			Description: "Description for Project Alpha",
		},
		{
			Name:        "Project Beta",
			Description: "Description for Project Beta",
		},
	}

	var wg sync.WaitGroup

	// Iterate through each project to insert it into the database
	for _, project := range projects {
		wg.Add(1)

		go func(project entities.Project) {
			defer wg.Done()

			// Insert project into the database
			_, err := db.Exec(`
				INSERT INTO projects (name, description, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4)
			`, project.Name, project.Description, time.Now(), time.Now())

			if err != nil {
				log.Printf("Failed to insert project %s: %v\n", project.Name, err)
				return
			}
			log.Printf("Successfully inserted project %s\n", project.Name)

			// Fetch the inserted project ID
			var projectId int
			err = db.QueryRow(`SELECT id FROM projects WHERE name = $1`, project.Name).Scan(&projectId)
			if err != nil {
				log.Printf("Failed to fetch project ID for %s: %v\n", project.Name, err)
				return
			}

			// Fetch all users (or filter based on criteria) for association with the project
			users := []entities.User{}
			rows, err := db.Query(`SELECT id FROM users`)
			if err != nil {
				log.Printf("Failed to fetch users: %v\n", err)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var user entities.User
				if err := rows.Scan(&user.ID); err != nil {
					log.Printf("Failed to scan user: %v\n", err)
					continue
				}
				users = append(users, user)
			}

			// Ensure no errors occurred during iteration over rows
			if err := rows.Err(); err != nil {
				log.Printf("Error iterating over users: %v\n", err)
				return
			}

			// Insert project-user associations
			for _, user := range users {
				// Insert project-user association into the users_projects table
				_, err = db.Exec(`
					INSERT INTO users_projects (project_id, user_id) 
					VALUES ($1, $2)
				`, projectId, user.ID)

				if err != nil {
					log.Printf("Failed to associate user %d with project %d: %v\n", user.ID, projectId, err)
					continue
				}
				log.Printf("Successfully associated user %d with project %d\n", user.ID, projectId)
			}
		}(project)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}
