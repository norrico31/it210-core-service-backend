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
