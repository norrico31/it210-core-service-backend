package seeders

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedTasks(db *sql.DB) error {
	// Fetch statuses
	statuses := make(map[string]int)
	statusRows, err := db.Query("SELECT id, name FROM statuses")
	if err != nil {
		log.Printf("Failed to fetch statuses: %v\n", err)
		return err
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var id int
		var name string
		if err := statusRows.Scan(&id, &name); err != nil {
			log.Printf("Failed to scan status: %v\n", err)
			return err
		}
		statuses[name] = id
	}

	// Fetch users
	users := map[string]int{}
	userRows, err := db.Query("SELECT id, firstName FROM users")
	if err != nil {
		log.Printf("Failed to fetch users: %v\n", err)
		return err
	}
	defer userRows.Close()

	for userRows.Next() {
		var id int
		var firstName string
		if err := userRows.Scan(&id, &firstName); err != nil {
			log.Printf("Failed to scan user: %v\n", err)
			return err
		}
		users[firstName] = id
	}

	// Fetch projects
	projects := make(map[string]int)
	projectRows, err := db.Query("SELECT id, name FROM projects")
	if err != nil {
		log.Printf("Failed to fetch projects: %v\n", err)
		return err
	}
	defer projectRows.Close()

	for projectRows.Next() {
		var id int
		var name string
		if err := projectRows.Scan(&id, &name); err != nil {
			log.Printf("Failed to scan project: %v\n", err)
			return err
		}
		projects[name] = id
	}

	user1 := users["Chester"]
	user2 := users["Mary Grace"]

	proj1 := projects["Project 124 Interpreter"]
	proj12 := projects["Project 210 Web App DOSTV"]

	// Sample task data
	tasks := []entities.Task{
		{
			Title:       "Design Database Schema",
			Description: "Design the database schema for the project.",
			StatusID:    statuses["In Progress"],
			UserID:      &user1,
			ProjectID:   proj12,
			SubTask: []entities.SubTask{
				{Title: "Define ER model", StatusID: statuses["Not Started"]},
				{Title: "Setup tables", StatusID: statuses["In Progress"]},
				{Title: "Define constraints", StatusID: statuses["Completed"]},
			},
		},
		{
			Title:       "Develop API Endpoints",
			Description: "Develop all required API endpoints.",
			StatusID:    statuses["In Progress"],
			UserID:      &user2,
			ProjectID:   proj1,
			SubTask: []entities.SubTask{
				{Title: "Setup router", StatusID: statuses["Not Started"]},
				{Title: "Create handlers", StatusID: statuses["In Progress"]},
				{Title: "Write tests", StatusID: statuses["Not Started"]},
			},
		},
	}

	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go func(task entities.Task) {
			defer wg.Done()

			// Insert task
			var taskID int
			err := db.QueryRow(`
				INSERT INTO tasks (title, description, statusId, userId, projectId, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
			`, task.Title, task.Description, task.StatusID, task.UserID, task.ProjectID, time.Now(), time.Now()).Scan(&taskID)
			if err != nil {
				log.Printf("Failed to insert task %s: %v\n", task.Title, err)
				return
			}
			log.Printf("Inserted task: %s with ID: %d\n", task.Title, taskID)

			// Insert subtasks
			for _, subTask := range task.SubTask {
				_, err := db.Exec(`
					INSERT INTO subtasks (taskId, statusId, title, createdAt, updatedAt)
					VALUES ($1, $2, $3, $4, $5)
				`, taskID, subTask.StatusID, subTask.Title, time.Now(), time.Now())
				if err != nil {
					log.Printf("Failed to insert subtask %s for task %s: %v\n", subTask.Title, task.Title, err)
				} else {
					log.Printf("Inserted subtask: %s for task: %s\n", subTask.Title, task.Title)
				}
			}
		}(task)
	}

	wg.Wait()
	return nil
}
