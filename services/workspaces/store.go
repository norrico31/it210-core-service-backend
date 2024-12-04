package workspaces

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/lib/pq"
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

func (s *Store) GetWorkspace(projectId int) ([]entities.Workspace, error) {
	// Step 1: Fetch Workspaces
	workspacesQuery := `
    SELECT 
        id AS workspace_id,
        name AS workspace_name,
        description AS workspace_description,
        projectId AS workspace_project_id,
        colOrder AS workspace_col_order,
        createdAt AS workspace_createdAt,
        updatedAt AS workspace_updatedAt,
        deletedAt AS workspace_deletedAt
    FROM workspaces
    WHERE projectId = $1 AND deletedAt IS NULL
    ORDER BY createdAt DESC;
`

	rows, err := s.db.Query(workspacesQuery, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspaces: %w", err)
	}
	defer rows.Close()

	var workspaceIDs []int
	workspacesMap := make(map[int]*entities.Workspace)

	for rows.Next() {
		workspace := entities.Workspace{}
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
			return nil, fmt.Errorf("failed to scan workspace data: %w", err)
		}
		workspacesMap[workspace.ID] = &workspace
		workspacesMap[workspace.ID].Tasks = []entities.Task{} // Initialize empty task list
		workspaceIDs = append(workspaceIDs, workspace.ID)
	}

	// Early return if no workspaces
	if len(workspaceIDs) == 0 {
		return []entities.Workspace{}, nil
	}

	// Step 2: Fetch Tasks (same as before)
	tasksQuery := `
    SELECT 
        id AS task_id,
        workspaceId AS task_workspace_id,
        title AS task_title,
        description AS task_description,
        userId AS task_user_id,
        priorityId AS task_priority_id,
        taskOrder AS task_order,
        createdAt AS task_createdAt,
        updatedAt AS task_updatedAt,
        deletedAt task_deletedAt
    FROM tasks
    WHERE workspaceId = ANY($1) AND deletedAt IS NULL;
`

	taskRows, err := s.db.Query(tasksQuery, pq.Array(workspaceIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer taskRows.Close()

	// Map tasks to workspaces
	for taskRows.Next() {
		task := entities.Task{}
		var workspaceID int
		err := taskRows.Scan(
			&task.ID,
			&workspaceID,
			&task.Title,
			&task.Description,
			&task.UserID,
			&task.PriorityID,
			&task.TaskOrder,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task data: %w", err)
		}
		if workspace, exists := workspacesMap[workspaceID]; exists {
			workspace.Tasks = append(workspace.Tasks, task)
		}
	}

	// Step 3: Sort tasks by taskOrder in ascending order
	for _, workspace := range workspacesMap {
		sort.Slice(workspace.Tasks, func(i, j int) bool {
			return workspace.Tasks[i].TaskOrder < workspace.Tasks[j].TaskOrder
		})
	}

	// Convert map to slice
	var workspaces []entities.Workspace
	for _, workspace := range workspacesMap {
		workspaces = append(workspaces, *workspace)
	}

	return workspaces, nil
}

func (s *Store) CreateWorkspace(payload entities.WorkspacePayload) (*entities.Workspace, error) {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Insert the workspace data into the workspaces table
	var workspaceID int
	err = tx.QueryRow(`
		INSERT INTO workspaces (name, description, projectId, colOrder) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`,
		payload.Name, payload.Description, payload.ProjectID, payload.ColOrder,
	).Scan(&workspaceID)

	if err != nil {
		// Rollback the transaction in case of an error
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return nil, fmt.Errorf("failed to insert workspace: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Fetch and return the newly created workspace
	workspace := &entities.Workspace{
		ID:          workspaceID,
		Name:        payload.Name,
		Description: payload.Description,
		ProjectID:   payload.ProjectID,
		ColOrder:    payload.ColOrder,
		CreatedAt:   time.Now(), // Assuming this field is not overwritten by the DB
		UpdatedAt:   time.Now(), // Assuming this field is not overwritten by the DB
	}

	return workspace, nil
}

// TODO DAPAT MATCH UNG PROJECTID SA WORKSPACE ID
func (s *Store) UpdateWorkspace(payload entities.WorkspacePayload) error {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute the update query
	_, err = tx.Exec(`
		UPDATE workspaces 
		SET name = $1, description = $2, colOrder = $3, updatedAt = CURRENT_TIMESTAMP 
		WHERE id = $4 AND projectId = $5 AND deletedAt IS NULL`,
		payload.Name, payload.Description, payload.ColOrder, payload.ID, payload.ProjectID,
	)
	if err != nil {
		// Rollback the transaction in case of an error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("update error: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
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

// func (s *Store) TaskDragNDrop(workspaceId, sourceTaskId, destinationTaskId int) error {
// 	// Step 1: Fetch tasks for the given workspace ordered by taskOrder
// 	tasksQuery := `
// 		SELECT
// 			id,
// 			taskOrder
// 		FROM tasks
// 		WHERE workspaceId = $1
// 		ORDER BY taskOrder;
// 	`
// 	rows, err := s.db.Query(tasksQuery, workspaceId)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch tasks for workspace %d: %w", workspaceId, err)
// 	}
// 	defer rows.Close()

// 	// Step 2: Create a map of taskID to taskOrder and a list of tasks in order
// 	var tasks []struct {
// 		ID        int
// 		TaskOrder int
// 	}
// 	for rows.Next() {
// 		var task struct {
// 			ID        int
// 			TaskOrder int
// 		}
// 		if err := rows.Scan(&task.ID, &task.TaskOrder); err != nil {
// 			return fmt.Errorf("failed to scan task data: %w", err)
// 		}
// 		tasks = append(tasks, task)
// 	}

// 	// Step 3: Identify the source and destination tasks
// 	var sourceIndex, destinationIndex int
// 	var sourceTaskOrder, destinationTaskOrder int
// 	for i, task := range tasks {
// 		if task.ID == sourceTaskId {
// 			sourceIndex = i
// 			sourceTaskOrder = task.TaskOrder
// 		}
// 		if task.ID == destinationTaskId {
// 			destinationIndex = i
// 			destinationTaskOrder = task.TaskOrder
// 		}
// 	}

// 	// If the source or destination task is not found, return an error
// 	if sourceTaskOrder == 0 || destinationTaskOrder == 0 {
// 		return fmt.Errorf("source or destination task not found")
// 	}

// 	// Step 4: Update taskOrder for tasks
// 	// If the source is being moved before or after the destination task
// 	if sourceIndex < destinationIndex {
// 		// Moving task after destination - shift tasks between source and destination
// 		for i := sourceIndex + 1; i <= destinationIndex; i++ {
// 			tasks[i].TaskOrder--
// 		}
// 	} else {
// 		// Moving task before destination - shift tasks between destination and source
// 		for i := destinationIndex; i < sourceIndex; i++ {
// 			tasks[i].TaskOrder++
// 		}
// 	}

// 	// Update the source task to the destination task's position
// 	tasks[sourceIndex].TaskOrder = destinationTaskOrder

// 	// Step 5: Update taskOrder in the database
// 	for _, task := range tasks {
// 		_, err := s.db.Exec(`UPDATE tasks SET taskOrder = $1 WHERE id = $2`, task.TaskOrder, task.ID)
// 		if err != nil {
// 			return fmt.Errorf("failed to update taskOrder for task %d: %w", task.ID, err)
// 		}
// 	}

// 	return nil
// }

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
