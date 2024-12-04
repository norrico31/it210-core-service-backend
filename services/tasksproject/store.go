package tasksproject

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

func (s *Store) GetTasksProject(projectId int) ([]*entities.TasksProject, error) {
	// SQL query without subtasks
	print("PWEDE NA KO MAG PRINT DITO?")
	query := fmt.Sprintf(`
        SELECT 
			t.id, t.title, t.description, t.userId, t.priorityId, t.projectId, t.createdAt, t.updatedAt, t.deletedAt, t.deletedBy
        FROM tasks t
    `)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasksMap := make(map[int]*entities.TasksProject)
	for rows.Next() {
		tasksProject := entities.TasksProject{}

		err := rows.Scan(
			&tasksProject.ID, &tasksProject.Name, &tasksProject.Description, &tasksProject.UserID, &tasksProject.PriorityID, &tasksProject.ProjectID, &tasksProject.CreatedAt,
			&tasksProject.UpdatedAt, &tasksProject.DeletedAt, &tasksProject.DeletedBy,
		)

		if err != nil {
			log.Printf("Failed to scan tasksProject: %v", err)
			continue
		}
		tasksMap[tasksProject.ID] = &tasksProject
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over tasksProject rows: %v", err)
	}

	tasksProjectList := make([]*entities.TasksProject, 0, len(tasksMap))
	for _, tasksProject := range tasksMap {
		tasksProjectList = append(tasksProjectList, tasksProject)
	}

	return tasksProjectList, nil
}

func (s *Store) GetTaskProject(projectId, taskId int) (*entities.TasksProject, error) {
	query := fmt.Sprintf(`
        SELECT 
			t.id, t.title, t.description, t.userId, t.priorityId, t.projectId, t.createdAt, t.updatedAt, t.deletedAt, t.deletedBy,

			u.id AS user_id, u.firstName, u.lastName, u.email, u.age, u.lastActiveAt, u.createdAt AS user_createdAt, u.updatedAt AS user_updatedAt, u.deletedAt AS user_deletedAt,

			p.id priority_id, p.name priority_name, p.description priority_description, p.createdAt priority_createdAt, p.updatedAt priority_updatedAt, p.deletedAt priority_deletedAt,

			w.id workspace_id, w.name workspace_name, w.description workspace_description

			FROM tasks t
			LEFT JOIN
				 users u ON u.id = t.userId
			LEFT JOIN
				priorities p ON p.id = t.priorityId
			LEFT JOIN
				workspaces w ON w.id = t.workspaceId
			WHERE t.id = $1 AND t.deletedAt IS NULL
    `)

	row := s.db.QueryRow(query, projectId)

	tasksProject := &entities.TasksProject{}
	var user entities.User
	var priority entities.Priority
	var workspace entities.Workspace

	err := row.Scan(
		&tasksProject.ID, &tasksProject.Name, &tasksProject.Description, &tasksProject.UserID, &tasksProject.PriorityID, &tasksProject.ProjectID, &tasksProject.CreatedAt,
		&tasksProject.UpdatedAt, &tasksProject.DeletedAt, &tasksProject.DeletedBy,

		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.LastActiveAt, &user.CreatedAt,
		&user.UpdatedAt, &user.DeletedAt,

		&priority.ID, &priority.Name, &priority.Description, &priority.CreatedAt, &priority.UpdatedAt, &priority.DeletedAt,

		&workspace.ID, &workspace.Name, &workspace.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tasksProject with ID %d not found", projectId)
		}
		return nil, fmt.Errorf("failed to retrieve tasksProject: %v", err)
	}

	if priority.ID != 0 {
		tasksProject.Priority = priority
	}
	if user.ID != 0 {
		tasksProject.User = user
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
		INSERT INTO tasks (title, description, userId, priorityId, workspaceId, taskOrder)
		VALUES ($1, $2, $3, $4, $5, NULL)
		RETURNING id, title, description, userId, priorityId, workspaceId, taskOrder, createdAt, updatedAt
	`
	err = tx.QueryRow(
		query,
		payload.Name,
		payload.Description,
		sql.NullInt64{Int64: int64(payload.UserID), Valid: payload.UserID != 0},
		payload.PriorityID,
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

	_, err = tx.Exec(`UPDATE tasks SET title = $1, description = $2, userId = $3, priorityId = $4, updatedAt = CURRENT_TIMESTAMP WHERE id = $5`,
		payload.Name,
		payload.Description,
		payload.UserID,
		payload.PriorityID,
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
	err = tx.QueryRow("UPDATE tasks SET deletedAt = NULL WHERE id = $1 RETURNING id, title, description, userId, createdAt, updatedAt, deletedAt", id).Scan(
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
