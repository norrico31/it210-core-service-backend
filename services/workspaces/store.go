package workspaces

import (
	"database/sql"
	"fmt"

	"github.com/norrico31/it210-core-service-backend/entities"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// TODO ADD THE FUNCTIONALITY OF DRAG N DROP HERE FOR COLUMN IN WORKSPACE
func (s *Store) GetWorkspaces() ([]entities.Workspace, error) {
	queryWorkspaces := `
        SELECT 
            w.id, w.name, w.description, w.projectId, w.colOrder,
            w.createdAt, w.updatedAt, w.deletedAt
        FROM workspaces w
        ORDER BY w.createdAt DESC
    `

	rows, err := s.db.Query(queryWorkspaces)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []entities.Workspace
	workspaceMap := make(map[int]*entities.Workspace)

	// Process workspaces
	for rows.Next() {
		var workspace entities.Workspace

		err := rows.Scan(
			&workspace.ID,
			&workspace.Name,
			&workspace.Description,
			&workspace.ProjectID,
			&workspace.ColOrder,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
			&workspace.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workspace: %w", err)
		}

		workspace.Tasks = []entities.Task{} // Initialize Tasks slice
		workspaces = append(workspaces, workspace)
		workspaceMap[workspace.ID] = &workspaces[len(workspaces)-1] // Map ID to workspace pointer
	}

	return workspaces, nil
}

// func (s *Store) GetWorkspaces(projectId int) ([]entities.Workspace, error) {
// 	println("ProjectId: %s", projectId)
// 	query := `
// 		SELECT
// 			w.id, w.name, w.description, w.projectId, w.createdAt, w.updatedAt,
// 			p.id AS project_id, p.name AS project_name, p.description AS project_description,
// 			p.progress AS project_progress, p.dateStarted AS project_date_started, p.dateDeadline AS project_date_deadline,
// 			t.id AS task_id, t.title AS task_title, t.description AS task_description,
// 			t.userId AS task_user_id, t.priorityId AS task_priority_id,
// 			t.taskOrder AS task_order, t.createdAt AS task_created_at
// 		FROM workspaces w
// 		LEFT JOIN projects p ON p.id = w.projectId
// 		LEFT JOIN tasks t ON t.workspaceId = w.id
// 		WHERE w.deletedAt IS NULL AND w.projectId = $1
// 		ORDER BY w.createdAt DESC
// 	`

// 	rows, err := s.db.Query(query, projectId)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query workspaces for project ID %d: %w", projectId, err)
// 	}
// 	defer rows.Close()

// 	workspaces := []entities.Workspace{}
// 	workspaceMap := make(map[int]*entities.Workspace)

// 	for rows.Next() {
// 		var workspace entities.Workspace
// 		var project entities.Project
// 		var task entities.Task
// 		var taskID sql.NullInt64

// 		err := rows.Scan(
// 			&workspace.ID, &workspace.Name, &workspace.Description, &workspace.ProjectID,
// 			&workspace.CreatedAt, &workspace.UpdatedAt,
// 			&project.ID, &project.Name, &project.Description, &project.Progress,
// 			&project.DateStarted, &project.DateDeadline,
// 			&taskID, &task.Title, &task.Description, &task.UserID,
// 			&task.PriorityID, &task.TaskOrder, &task.CreatedAt,
// 		)
// 		if err != nil {
// 			log.Printf("Failed to scan workspace or task: %v", err)
// 			continue
// 		}

// 		// Ensure workspace is only added once
// 		if _, exists := workspaceMap[workspace.ID]; !exists {
// 			workspace.Project = project
// 			workspace.Tasks = []entities.Task{}
// 			workspaceMap[workspace.ID] = &workspace
// 			workspaces = append(workspaces, workspace)
// 		}

// 		// Add task if it has a valid ID
// 		if taskID.Valid {
// 			task.ID = int(taskID.Int64)
// 			workspaceMap[workspace.ID].Tasks = append(workspaceMap[workspace.ID].Tasks, task)
// 		}
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to iterate over workspace rows: %w", err)
// 	}

// 	return workspaces, nil
// }

func (s *Store) GetWorkspace(projectId int) (*entities.Workspace, error) {
	// Query to get workspaces
	queryWorkspaces := `
        SELECT 
            w.id, w.name, w.description, w.projectId, w.colOrder,
            w.createdAt, w.updatedAt, w.deletedAt
        FROM workspaces w
        WHERE w.projectId = $1 AND w.deletedAt IS NULL
        ORDER BY w.createdAt DESC
    `

	rows, err := s.db.Query(queryWorkspaces, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspaces for project ID %d: %w", projectId, err)
	}
	defer rows.Close()

	var workspaces []entities.Workspace
	workspaceMap := make(map[int]*entities.Workspace)

	// Process workspaces
	for rows.Next() {
		var workspace entities.Workspace

		err := rows.Scan(
			&workspace.ID,
			&workspace.Name,
			&workspace.Description,
			&workspace.ProjectID,
			&workspace.ColOrder,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
			&workspace.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workspace: %w", err)
		}

		workspace.Tasks = []entities.Task{} // Initialize Tasks slice
		workspaces = append(workspaces, workspace)
		workspaceMap[workspace.ID] = &workspaces[len(workspaces)-1] // Map ID to workspace pointer
	}

	return nil, nil
}

func (s *Store) CreateWorkspace(payload entities.WorkspacePayload) (*entities.Workspace, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO workspaces (name, description) VALUES ($1, $2)", payload.Name, payload.Description)
	if err != nil {
		// If there's an error, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return nil, err
	}

	// Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return &entities.Workspace{Name: payload.Name, Description: payload.Description}, err
	}

	return &entities.Workspace{Name: payload.Name, Description: payload.Description}, err
}

func (s *Store) UpdateWorkspace(payload entities.WorkspacePayload) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE workspaces SET name = $1, description = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3", payload.Name, payload.Description, payload.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (s *Store) DeleteWorkspace(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE workspaces SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1", id)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("delete error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (s *Store) RestoreWorkspace(id int) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE workspaces SET deletedAt = NULL WHERE id = $1", id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("delete error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	// Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func scanRowIntoWorkspace(rows *sql.Rows, workspace *entities.Workspace) error {
	return rows.Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Description,
		&workspace.ProjectID,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
}
