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
	// SQL query without subtasks
	query := fmt.Sprintf(`
        SELECT 
			t.id, t.title, t.description, t.userId, t.priorityId, t.workspaceId, t.taskOrder, t.createdAt, t.updatedAt, t.deletedAt, t.deletedBy
        FROM tasks t
    `)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasksMap := make(map[int]*entities.Task)
	for rows.Next() {
		task := entities.Task{}

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.UserID, &task.PriorityID, &task.WorkspaceID, &task.TaskOrder, &task.CreatedAt,
			&task.UpdatedAt, &task.DeletedAt, &task.DeletedBy,
		)

		if err != nil {
			log.Printf("Failed to scan task: %v", err)
			continue
		}
		tasksMap[task.ID] = &task
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over tasks rows: %v", err)
	}

	tasks := make([]*entities.Task, 0, len(tasksMap))
	for _, task := range tasksMap {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Store) GetTask(id int) (*entities.Task, error) {
	query := fmt.Sprintf(`
        SELECT 
			t.id, t.title, t.description, t.userId, t.priorityId, t.workspaceId, t.taskOrder, t.createdAt, t.updatedAt, t.deletedAt, t.deletedBy,

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

	row := s.db.QueryRow(query, id)

	task := &entities.Task{}
	var user entities.User
	var priority entities.Priority
	var workspace entities.Workspace

	err := row.Scan(
		&task.ID, &task.Title, &task.Description, &task.UserID, &task.PriorityID, &task.WorkspaceID, &task.TaskOrder, &task.CreatedAt,
		&task.UpdatedAt, &task.DeletedAt, &task.DeletedBy,

		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.LastActiveAt, &user.CreatedAt,
		&user.UpdatedAt, &user.DeletedAt,

		&priority.ID, &priority.Name, &priority.Description, &priority.CreatedAt, &priority.UpdatedAt, &priority.DeletedAt,

		&workspace.ID, &workspace.Name, &workspace.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve task: %v", err)
	}

	if priority.ID != 0 {
		task.Priority = priority
	}
	if user.ID != 0 {
		task.User = user
	}
	if workspace.ID != 0 {
		task.Workspace = workspace
	}
	return task, nil
}

func (s *Store) TaskCreate(payload entities.TaskCreatePayload) (*entities.Task, error) {
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

	task := entities.Task{}
	query := `
		INSERT INTO tasks (title, description, userId, priorityId, workspaceId, taskOrder)
		VALUES ($1, $2, $3, $4, $5, NULL)
		RETURNING id, title, description, userId, priorityId, workspaceId, taskOrder, createdAt, updatedAt
	`
	err = tx.QueryRow(
		query,
		payload.Title,
		payload.Description,
		sql.NullInt64{Int64: int64(payload.UserID), Valid: payload.UserID != 0},
		payload.PriorityID,
		payload.WorkspaceID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.UserID,
		&task.PriorityID,
		&task.WorkspaceID,
		&task.TaskOrder,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &task, nil
}

func (s *Store) TaskUpdate(payload entities.TaskUpdatePayload) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	fmt.Println("payload: ", payload.ID)
	fmt.Println("payload: ", payload.PriorityID)
	fmt.Println("payload: ", payload.UserID)
	fmt.Println("payload: ", payload.WorkspaceID)
	_, err = tx.Exec(`UPDATE tasks SET title = $1, description = $2, userId = $3, priorityId = $4, workspaceId = $5, updatedAt = CURRENT_TIMESTAMP WHERE id = $6`,
		payload.Title,
		payload.Description,
		payload.UserID,
		payload.PriorityID,
		payload.WorkspaceID,
		payload.ID,
	)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Store) TaskDelete(id int) (*entities.Task, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	task := &entities.Task{}
	err = tx.QueryRow("UPDATE tasks SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id, title, description, userId, createdAt, updatedAt, deletedAt", id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.UserID,
		// &task.ProjectID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
	)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error deleting : %v, rollback error: %v", err, rollbackErr)
		}

		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Store) TaskRestore(id int) (*entities.Task, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	task := entities.Task{}
	err = tx.QueryRow("UPDATE tasks SET deletedAt = NULL WHERE id = $1 RETURNING id, title, description, userId, createdAt, updatedAt, deletedAt", id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.UserID,
		// &task.ProjectID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error restoring: %v rollback error: %v", err, rollbackErr)
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return &task, nil
	}

	return &task, nil
}

func scanRowIntoTask(rows *sql.Rows, task *entities.Task) error {
	return rows.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.UserID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
		&task.DeletedBy,
	)
}
