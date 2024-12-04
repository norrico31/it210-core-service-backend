package tasksproject

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetTasksProject(projectId int) ([]*entities.TasksProject, error) {
	// SQL query to get tasks based on projectId and include user and priority details as nested objects
	query := fmt.Sprintf(`
        SELECT 
			pt.id task_id, 
			pt.name task_name, 
			pt.description task_description, 
			pt.userId task_user_id, 
			pt.priorityId task_priority_id, 
			pt.projectId task_project_id, 
			pt.createdAt task_createdAt, 
			pt.updatedAt task_updatedAt, 
			pt.deletedAt task_deletedAt,
			u.firstName user_firstname,
			u.lastName user_lastname,
			u.age user_age,
			u.email user_email,
			p.name priority_name,
			p.description priority_description
        FROM project_tasks pt
        JOIN users u ON pt.userId = u.id
        JOIN priorities p ON pt.priorityId = p.id
        WHERE pt.projectId = $1 AND pt.deletedAt IS NULL AND u.deletedAt IS NULL AND p.deletedAt IS NULL
    `)

	rows, err := s.db.Query(query, projectId) // Pass projectId as a parameter
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasksProjectList := []*entities.TasksProject{}
	for rows.Next() {
		tasksProject := entities.TasksProject{}
		var userFirstName, userLastName, userEmail string
		var userAge int
		var priorityName, priorityDescription string

		err := rows.Scan(
			&tasksProject.ID, &tasksProject.Name, &tasksProject.Description, &tasksProject.UserID,
			&tasksProject.PriorityID, &tasksProject.ProjectID, &tasksProject.CreatedAt,
			&tasksProject.UpdatedAt, &tasksProject.DeletedAt,
			&userFirstName, &userLastName, &userAge, &userEmail,
			&priorityName, &priorityDescription,
		)

		if err != nil {
			log.Printf("Failed to scan tasksProject: %v", err)
			continue
		}

		// Set user info in a nested User object
		tasksProject.User = entities.User{
			ID:        *tasksProject.UserID,
			FirstName: userFirstName,
			LastName:  userLastName,
			Age:       userAge,
			Email:     userEmail,
		}

		// Set priority info in a nested Priority object
		tasksProject.Priority = entities.Priority{
			ID:          tasksProject.PriorityID,
			Name:        priorityName,
			Description: priorityDescription,
		}

		tasksProjectList = append(tasksProjectList, &tasksProject)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over tasksProject rows: %v", err)
	}

	return tasksProjectList, nil
}

func (s *Store) GetTaskProject(taskId int) (*entities.TasksProject, error) {
	// SQL query to get a single task with related user, priority, and project details
	query := fmt.Sprintf(`
        SELECT 
			pt.id, pt.name, pt.description, pt.userId, pt.priorityId, pt.projectId, pt.createdAt, pt.updatedAt, pt.deletedAt, pt.deletedBy,
			u.firstName user_firstname, u.lastName user_lastname, u.age user_age, u.email user_email,
			p.name priority_name, p.description priority_description,
			pr.name project_name, pr.description project_description, pr.progress project_progress, pr.url project_url, pr.dateStarted project_dateStarted, pr.dateDeadline project_dateDeadline
		FROM project_tasks pt
		LEFT JOIN users u ON u.id = pt.userId
		LEFT JOIN priorities p ON p.id = pt.priorityId
		LEFT JOIN projects pr ON pr.id = pt.projectId
		WHERE pt.id = $1 AND pt.deletedAt IS NULL
    `)

	row := s.db.QueryRow(query, taskId)

	tasksProject := &entities.TasksProject{}
	var userFirstName, userLastName, userEmail string
	var userAge int
	var priorityName, priorityDescription string
	var projectName, projectDescription, projectURL *string
	var projectProgress *float64
	var projectDateStarted, projectDateDeadline *time.Time

	// Scan the result into the TasksProject and related User, Priority, and Project fields
	err := row.Scan(
		&tasksProject.ID, &tasksProject.Name, &tasksProject.Description, &tasksProject.UserID, &tasksProject.PriorityID, &tasksProject.ProjectID, &tasksProject.CreatedAt,
		&tasksProject.UpdatedAt, &tasksProject.DeletedAt, &tasksProject.DeletedBy,

		&userFirstName, &userLastName, &userAge, &userEmail,
		&priorityName, &priorityDescription,

		&projectName, &projectDescription, &projectProgress, &projectURL, &projectDateStarted, &projectDateDeadline,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with ID %d not found", taskId)
		}
		return nil, fmt.Errorf("failed to retrieve task: %v", err)
	}

	// Set User info in a nested User object if user exists
	if userFirstName != "" && userLastName != "" {
		tasksProject.User = entities.User{
			ID:        *tasksProject.UserID,
			FirstName: userFirstName,
			LastName:  userLastName,
			Age:       userAge,
			Email:     userEmail,
		}
	}

	// Set Priority info in a nested Priority object if priority exists
	if priorityName != "" && priorityDescription != "" {
		tasksProject.Priority = entities.Priority{
			ID:          tasksProject.PriorityID,
			Name:        priorityName,
			Description: priorityDescription,
		}
	}

	// Set Project info in a nested Project object if project exists
	if projectName != nil && projectDescription != nil {
		tasksProject.Project = entities.Project{
			ID:           tasksProject.PriorityID,
			Name:         *projectName,
			Description:  *projectDescription,
			Progress:     projectProgress,
			Url:          projectURL,
			DateStarted:  projectDateStarted,
			DateDeadline: projectDateDeadline,
		}
	}

	return tasksProject, nil
}

func (s *Store) TasksProjectCreate(payload entities.TasksProjectCreatePayload) (*entities.TasksProject, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	tasksProject := entities.TasksProject{}
	query := `
		INSERT INTO project_tasks (name, description, userId, priorityId, projectId, createdAt, updatedAt)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, description, userId, priorityId, createdAt, updatedAt
	`
	err = tx.QueryRow(
		query,
		payload.Name,
		payload.Description,
		payload.UserID,
		payload.PriorityID,
		payload.ProjectID,
		time.Now(),
		time.Now(),
	).Scan(
		&tasksProject.ID,
		&tasksProject.Name,
		&tasksProject.Description,
		&tasksProject.UserID,
		&tasksProject.PriorityID,
		&tasksProject.CreatedAt,
		&tasksProject.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tasksProject: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &tasksProject, nil
}

func (s *Store) TasksProjectUpdate(payload entities.TasksProjectUpdatePayload) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE project_tasks 
		SET name = $1, description = $2, userId = $3, priorityId = $4, projectId = $5, updatedAt = CURRENT_TIMESTAMP 
		WHERE id = $6
		`,
		payload.Name,
		payload.Description,
		payload.UserID,
		payload.PriorityID,
		payload.ProjectID,
		payload.ID,
	)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("update error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Store) TasksProjectDelete(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE tasks SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1", id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error deleting tasksProject: %v, rollback error: %v", err, rollbackErr)
		}
		return nil
	}

	if err = tx.Commit(); err != nil {
		return nil
	}

	return nil
}

func (s *Store) TasksProjectRestore(id int) (*entities.TasksProject, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	tasksProject := entities.TasksProject{}
	err = tx.QueryRow("UPDATE tasks SET deletedAt = NULL WHERE id = $1 RETURNING id, name, description, userId, createdAt, updatedAt, deletedAt", id).Scan(
		&tasksProject.ID,
		&tasksProject.Name,
		&tasksProject.Description,
		&tasksProject.UserID,
		&tasksProject.CreatedAt,
		&tasksProject.UpdatedAt,
		&tasksProject.DeletedAt,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error restoring tasksProject: %v rollback error: %v", err, rollbackErr)
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return &tasksProject, nil
	}

	return &tasksProject, nil
}

func scanRowIntoTasksProject(rows *sql.Rows, tasksProject *entities.TasksProject) error {
	return rows.Scan(
		&tasksProject.ID,
		&tasksProject.Name,
		&tasksProject.Description,
		&tasksProject.UserID,
		&tasksProject.CreatedAt,
		&tasksProject.UpdatedAt,
		&tasksProject.DeletedAt,
		&tasksProject.DeletedBy,
	)
}
