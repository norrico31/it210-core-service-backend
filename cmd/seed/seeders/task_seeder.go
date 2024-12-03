package seeders

import (
	"database/sql"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedTasks(db *sql.DB) error {
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

	priorities := make(map[string]int)
	priorityRows, err := db.Query("SELECT id, name FROM priorities")
	if err != nil {
		log.Printf("Failed to fetch priorities: %v\n", err)
		return err
	}
	defer priorityRows.Close()

	for priorityRows.Next() {
		var id int
		var name string
		if err := priorityRows.Scan(&id, &name); err != nil {
			log.Printf("Failed to scan priorities: %v\n", err)
			return err
		}
		priorities[name] = id
	}

	workspaces := map[string]int{}
	workspaceRows, err := db.Query("SELECT id, name FROM workspaces")
	if err != nil {
		log.Printf("Failed to fetch workspaces: %v\n", err)
		return err
	}
	defer workspaceRows.Close()

	for workspaceRows.Next() {
		var id int
		var name string
		if err := workspaceRows.Scan(&id, &name); err != nil {
			log.Printf("Failed to scan workspaces: %v\n", err)
			return err
		}
		workspaces[name] = id
	}

	user1 := users["Chester"]
	user2 := users["Mary Grace"]

	high := priorities["High"]
	low := priorities["Low"]

	pending := workspaces["PENDING"]
	ongoing := workspaces["ONGOING"]

	// proj124 := projects["Project 124 Interpreter"]
	// proj210 := projects["Project 210 Web App DOSTV"]

	tasks := []entities.Task{
		{
			Title:       "Design Database Schema",
			Description: "Design the database schema for the project.",
			WorkspaceID: pending,
			UserID:      &user1,
			// ProjectID:   proj124,
			PriorityID: high,
		},
		{
			Title:       "Develop API Endpoints",
			Description: "Develop all required API endpoints.",
			WorkspaceID: ongoing,
			UserID:      &user2,
			// ProjectID:   proj210,
			PriorityID: low,
		},
		{
			Title:       "Task 3",
			Description: "Design the database schema for the project.",
			WorkspaceID: pending,
			UserID:      &user1,
			// ProjectID:   proj124,
			PriorityID: high,
		},
		{
			Title:       "Task 4",
			Description: "Develop all required API endpoints.",
			WorkspaceID: ongoing,
			UserID:      &user2,
			// ProjectID:   proj210,
			PriorityID: low,
		},
	}

	for _, task := range tasks {
		var taskID int
		err := db.QueryRow(`
			INSERT INTO tasks (title, description, workspaceId, userId, priorityId, createdAt, updatedAt)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
		`, task.Title, task.Description, task.WorkspaceID, task.UserID, task.PriorityID, time.Now(), time.Now()).Scan(&taskID)
		if err != nil {
			log.Printf("Failed to insert task %s: %v\n", task.Title, err)
			continue
		}
		log.Printf("Inserted task: %s \n", task.Title)
	}

	return nil
}
