package tasks

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/norrico31/it210-core-service-backend/entities"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetTasks() ([]*entities.Task, error) {
	query := fmt.Sprintf(`
		SELECT * FROM tasks;
	`)

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []*entities.Task{}
	for rows.Next() {
		task := entities.Task{}
		err := scanRowIntoTask(rows, &task)

		if err != nil {
			log.Printf("ailed to scan task: %v", err)
			continue
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over tasks rows: %v", err)
	}

	return tasks, nil
}

func (s *Store) GetTask(id int) (*entities.Task, error) {
	query := fmt.Sprintf(`
		SELECT * FROM tasks WHERE id = %v
	`, id)
	row := s.db.QueryRow(query)

	task := &entities.Task{}

	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.SubTask,
		&task.Description,
		&task.StatusID,
		&task.UserID,
		&task.ProjectID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("task ID not found")
	}

	if task.ID == 0 {
		return nil, fmt.Errorf("task ID not found ")
	}

	return task, nil
}

func (s *Store) TaskCreate(payload entities.TaskCreatePayload) (*entities.Task, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	// If the userId is not null (not 0), ensure it exists in the users table.
	if payload.UserID != 0 {
		var count int
		err := tx.QueryRow("SELECT COUNT(1) FROM users WHERE id = $1", payload.UserID).Scan(&count)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error checking user existence: %v", err)
		}
		if count == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("user with ID %d does not exist", payload.UserID)
		}
	}

	// Insert the new task into the database.
	task := entities.Task{}
	err = tx.QueryRow("INSERT INTO tasks (title, subTask, description, statusId, userId, projectId) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, subTask, description, statusId, userId, projectId, createdAt, updatedAt",
		payload.Title,
		payload.SubTask,
		payload.Description,
		payload.StatusID,
		// If the userId is 0 (i.e., NULL), pass NULL to the query
		func() interface{} {
			if payload.UserID == 0 {
				return nil // Allow NULL for userId
			}
			return payload.UserID
		}(),
		payload.ProjectID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.SubTask,
		&task.Description,
		&task.StatusID,
		&task.UserID, // Now this will be a pointer to int
		&task.ProjectID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("insert error: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit error: %v", err)
	}

	return &task, nil
}

// TODO: UPDATE TASK (HOW TO UPDATE SUBTASK(ARRAY OF STRINGS))

func scanRowIntoTask(rows *sql.Rows, task *entities.Task) error {
	return rows.Scan(
		&task.ID,
		&task.Title,
		&task.SubTask,
		&task.Description,
		&task.StatusID,
		&task.UserID,
		&task.ProjectID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
	)
}
