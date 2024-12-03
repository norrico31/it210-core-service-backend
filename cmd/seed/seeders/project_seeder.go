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
	if err != nil {
		log.Printf("Failed to fetch statuses: %v", err)
		return err
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var id int
		var name string

		if err = statusRows.Scan(&id, &name); err != nil {
			log.Printf("Failed to scan statuses: %v", err)
			return err
		}
		statuses[name] = id
	}

	active := statuses["Active"]
	notStarted := statuses["Not Started"]

	segments := make(map[string]int)
	segmentRows, err := db.Query(`
		SELECT
			id, name
		FROM segments
	`)
	if err != nil {
		log.Printf("Failed to fetch segments: %v", err)
		return err
	}
	defer segmentRows.Close()

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
			SegmentID:   bantayBulkan,
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
	for _, project := range projects {
		wg.Add(1)

		go func(project entities.Project) {
			defer wg.Done()

			_, err := db.Exec(`
				INSERT INTO projects (name, description, progress, statusId, segmentId, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, project.Name, project.Description, project.Progress, project.StatusID, project.SegmentID, time.Now(), time.Now())

			if err != nil {
				log.Printf("Failed to insert project %s: %v\n", project.Name, err)
				return
			}
			log.Printf("Successfully inserted project %s\n", project.Name)

			var projectId int
			err = db.QueryRow(`SELECT id FROM projects WHERE name = $1`, project.Name).Scan(&projectId)
			if err != nil {
				log.Printf("Failed to fetch project ID for %s: %v\n", project.Name, err)
				return
			}

			users := []int{}
			rows, err := db.Query(`SELECT id FROM users`)
			if err != nil {
				log.Printf("Failed to fetch users: %v\n", err)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var id int
				if err := rows.Scan(&id); err != nil {
					log.Printf("Failed to scan user: %v\n", err)
					return
				}
				users = append(users, id)
			}

			if err := rows.Err(); err != nil {
				log.Printf("Error iterating over users: %v\n", err)
				return
			}

			// Insert project-user associations
			for _, userId := range users {
				_, err = db.Exec(`
				INSERT INTO users_projects (project_id, user_id) 
				VALUES ($1, $2)
				`, projectId, userId)

				if err != nil {
					log.Printf("Failed to associate user %d with project %d: %v\n", userId, projectId, err)
					continue
				}
				log.Printf("Successfully associated user %d with project %d\n", userId, projectId)
			}
		}(project)
	}

	wg.Wait()

	return nil
}

func float64Ptr(value float64) *float64 {
	return &value
}
