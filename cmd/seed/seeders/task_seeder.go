package seeders

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedTasks(db *sql.DB) error {
	userIds := make(map[string]int)

	rows, err := db.Query("SELECT id, firstName FROM users")
	if err != nil {
		log.Printf("Failed to fetch users: %v\n", err)
		return err
	}
	defer rows.Close()

	// Map user firstNames to user IDs
	for rows.Next() {
		var id int
		var firstName string
		if err := rows.Scan(&id, &firstName); err != nil {
			log.Printf("Failed to scan user: %v\n", err)
			return err
		}
		userIds[firstName] = id
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed to iterate users: %v\n", err)
		return err
	}

	// Sample data for seeding
	tasks := []entities.Task{
		{
			Title:       "Design Database Schema",
			SubTask:     []string{"Define ER model", "Setup tables", "Define constraints"},
			Description: "desc 1",
			StatusID:    1,
			UserId:      userIds["Chester"], // change this for prod
			Projects:    []entities.Project{{ID: 1}},
		},
		{
			Title:       "Develop API Endpoints",
			SubTask:     []string{"Setup router", "Create handlers", "Write tests"},
			Description: "desc 2",
			StatusID:    2,
			UserId:      userIds["Mary Grace"], // change this for prod
			Projects:    []entities.Project{{ID: 1}},
		},
	}

	// Insert each task into the database
	for _, task := range tasks {
		// Convert subtask array to PostgreSQL array format
		subTaskArray := "{" + fmt.Sprintf("'%s'", task.SubTask[0])
		for _, sub := range task.SubTask[1:] {
			subTaskArray += fmt.Sprintf(", '%s'", sub)
		}
		subTaskArray += "}"

		// Ensure projectId is set correctly as it's a NOT NULL field
		if len(task.Projects) == 0 || task.Projects[0].ID == 0 {
			log.Printf("Task %s has no valid projectId set. Skipping task.\n", task.Title)
			continue
		}

		// Insert task into the database with projectId
		_, err := db.Exec(`
            INSERT INTO tasks (title, subTask, description, statusId, userId, projectId, createdAt, updatedAt)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        `, task.Title, subTaskArray, task.Description, task.StatusID, task.UserId, task.Projects[0].ID, time.Now(), time.Now())

		if err != nil {
			log.Printf("Failed to insert task %s: %v\n", task.Title, err)
			return err
		}
		log.Printf("Successfully inserted task %s\n", task.Title)

		// Fetch the inserted task ID to associate it with projects
		var taskId int
		err = db.QueryRow(`SELECT id FROM tasks WHERE title = $1`, task.Title).Scan(&taskId)
		if err != nil {
			log.Printf("Failed to fetch task ID for %s: %v\n", task.Title, err)
			return err
		}
		// Insert project-task associations (if the Projects slice is populated)
		for _, project := range task.Projects {
			_, err := db.Exec(`
                INSERT INTO project_tasks (project_id, task_id)
                VALUES ($1, $2)
            `, project.ID, taskId)

			if err != nil {
				log.Printf("Failed to associate task %d with project %d: %v\n", taskId, project.ID, err)
				return err
			}
			log.Printf("Successfully associated task %d with project %d\n", taskId, project.ID)
		}
	}

	return nil
}
