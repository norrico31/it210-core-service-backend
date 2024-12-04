package seeders

import (
	"database/sql"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedProjectTasks(db *sql.DB) error {
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

	user1 := users["Chester"]
	user2 := users["Mary Grace"]

	high := priorities["High"]
	low := priorities["Low"]

	proj124 := projects["Project 124 Interpreter"]
	proj210 := projects["Project 210 Web App DOSTV"]

	tasks := []entities.TasksProject{
		{
			Name:        "Task 1 for project interpreter",
			Description: "description for task 1",
			UserID:      &user1,
			ProjectID:   proj124,
			PriorityID:  high,
		},
		{
			Name:        "Task 1 for IT210 web app DOSTV",
			Description: "Develop all required API endpoints.",
			UserID:      &user2,
			ProjectID:   proj210,
			PriorityID:  low,
		},
		{
			Name:        "Task 2 for project interpreter",
			Description: "description for task 2 of interpreter",
			UserID:      &user1,
			ProjectID:   proj124,
			PriorityID:  high,
		},
		{
			Name:        "Task 2 for IT210 web app DOSTV",
			Description: "task 2 for project 210",
			UserID:      &user2,
			ProjectID:   proj210,
			PriorityID:  low,
		},
	}

	for _, task := range tasks {
		var taskID int
		err := db.QueryRow(`
			INSERT INTO project_tasks (name, description, userId, projectId, priorityId, createdAt, updatedAt)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
		`, task.Name, task.Description, task.UserID, task.ProjectID, task.PriorityID, time.Now(), time.Now()).Scan(&taskID)
		if err != nil {
			log.Printf("Failed to insert project task %s: %v\n", task.Name, err)
			continue
		}
		log.Printf("Inserted task: %s \n", task.Name)
	}

	return nil
}
