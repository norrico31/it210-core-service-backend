package seeders

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedProjects(db *sql.DB) error {
	statuses := make(map[string]int)
	statusRows, err := db.Query(`
		SELECT
			id, name
		FROM statuses	
	`)
	for statusRows.Next() {
		var id int
		var name string

		if err = statusRows.Scan(&id, &name); err != nil {
			log.Printf("Failed to scan statuses: %v", err)
			return err
		}
		statuses[name] = id
	}
	statusRows.Close()

	active := statuses["Active"]
	notStarted := statuses["Not Started"]

	segments := make(map[string]int)
	segmentRows, err := db.Query(`
	SELECT
	id, name
	FROM segments
	`)
	for segmentRows.Next() {
		var id int
		var name string
		if err = segmentRows.Scan(&id, &name); err != nil {
			log.Printf("Failed to scan segments: %v", err)
			return err
		}
		segments[name] = id
	}

	siyensikat := segments["SIYENSIKAT"]
	bantayBulkan := segments["BANTAY BULKAN"]

	segmentRows.Close()

	projects := []entities.Project{
		{
			Name:        "Project 124 Interpreter",
			Description: "Description for Project 124 nakakaiyak",
			Progress:    float64Ptr(31),
			StatusID:    notStarted,
			SegmentID:   siyensikat,
		},
		{
			Name:        "Project 210 Web App DOSTV",
			Description: "Description for web app dostv chill lang",
			Progress:    float64Ptr(20),
			StatusID:    active,
			SegmentID:   bantayBulkan,
		},
		{
			Name:        "Single Sleeping Barber",
			Description: "Description for single sleeping barber problem",
			Progress:    float64Ptr(18),
			StatusID:    notStarted,
			SegmentID:   siyensikat,
		},
		{
			Name:        "CMSC 124 Interpreter Project",
			Description: "Description for interpreter",
			Progress:    float64Ptr(10),
			StatusID:    active,
		},
		{
			Name:        "CMSC 124 Messenger APP Erlang",
			Description: "Description for messenger app in erlang",
			Progress:    float64Ptr(18),
			StatusID:    notStarted,
			SegmentID:   bantayBulkan,
		},
		{
			Name:        "CMSC 124 Rust Superior",
			Description: "Description for rustlings",
			Progress:    float64Ptr(50),
			StatusID:    active,
			SegmentID:   siyensikat,
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
				INSERT INTO projects (name, description, progress, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4, $5)
			`, project.Name, project.Description, project.Progress, time.Now(), time.Now())

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

func float64Ptr(value float64) *float64 {
	return &value
}
