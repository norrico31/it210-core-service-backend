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

// func (s *Store) GetTasks(str string) ([]*entities.Task, error) {
// 	// SQL query with joins to fetch associated data
// 	query := fmt.Sprintf(`
//         SELECT t.id, t.title, t.description, t.statusId, t.userId, t.projectId, t.createdAt, t.updatedAt, t.deletedAt, t.deletedBy,
//                st.id AS subtask_id, st.taskId, st.title AS subtask_title, st.statusId AS subtask_statusId, st.createdAt AS subtask_createdAt, st.updatedAt AS subtask_updatedAt, st.deletedAt AS subtask_deletedAt, st.deletedBy AS subtask_deletedBy,
//                s.id AS status_id, s.name AS status_name, s.description AS status_description, s.createdAt AS status_createdAt, s.updatedAt AS status_updatedAt, s.deletedAt AS status_deletedAt,
//                u.id AS user_id, u.firstName, u.lastName, u.email, u.age, u.lastActiveAt, u.createdAt AS user_createdAt, u.updatedAt AS user_updatedAt, u.deletedAt AS user_deletedAt,
//                p.id AS project_id, p.name AS project_name, p.description AS project_description, p.createdAt AS project_createdAt, p.updatedAt AS project_updatedAt, p.deletedAt AS project_deletedAt
//         FROM tasks t
//         LEFT JOIN subtasks st ON st.taskId = t.id
//         LEFT JOIN statuses s ON s.id = t.statusId
//         LEFT JOIN users u ON u.id = t.userId
//         LEFT JOIN projects p ON p.id = t.projectId
//         WHERE t.deletedAt %v
//     `, str)

// 	rows, err := s.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	tasksMap := make(map[int]*entities.Task) // Using a map to avoid duplication of tasks
// 	for rows.Next() {
// 		task := entities.Task{}

// 		// Create variables to hold all the values we will scan into
// 		var subtask entities.SubTask
// 		var status entities.Status
// 		var user entities.User
// 		var project entities.Project

// 		// Scan the result set into the task and associated entities
// 		err := rows.Scan(
// 			&task.ID, &task.Title, &task.Description, &task.StatusID, &task.UserID, &task.ProjectID, &task.CreatedAt,
// 			&task.UpdatedAt, &task.DeletedAt, &task.DeletedBy,
// 			// Subtask
// 			&subtask.ID, &subtask.TaskID, &subtask.Title, &subtask.StatusID, &subtask.CreatedAt, &subtask.UpdatedAt,
// 			&subtask.DeletedAt, &subtask.DeletedBy,
// 			// Status
// 			&status.ID, &status.Name, &status.Description, &status.CreatedAt, &status.UpdatedAt, &status.DeletedAt,
// 			// User
// 			&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.LastActiveAt, &user.CreatedAt,
// 			&user.UpdatedAt, &user.DeletedAt,
// 			// Project
// 			&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt,
// 		)

// 		if err != nil {
// 			log.Printf("Failed to scan task: %v", err)
// 			continue
// 		}

// 		// Check if the task already exists in the map (to avoid duplicates)
// 		if existingTask, exists := tasksMap[task.ID]; exists {
// 			// Append subtask only if it's not already added
// 			existingTask.SubTask = append(existingTask.SubTask, subtask)
// 		} else {
// 			// Assign the related entities for new tasks
// 			task.SubTask = append(task.SubTask, subtask)
// 			task.Status = status
// 			task.User = user
// 			task.Project = project
// 			// Add the task to the map
// 			tasksMap[task.ID] = &task
// 		}
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to iterate over tasks rows: %v", err)
// 	}

// 	// Convert the map to a slice
// 	tasks := make([]*entities.Task, 0, len(tasksMap))
// 	for _, task := range tasksMap {
// 		tasks = append(tasks, task)
// 	}

// 	return tasks, nil
// }

func (s *Store) GetTasks(str string) ([]*entities.Task, error) {
	// SQL query without subtasks
	query := fmt.Sprintf(`
        SELECT t.id, t.title, t.description, t.statusId, t.userId, t.projectId, t.createdAt, t.updatedAt, t.deletedAt, t.deletedBy,
               s.id AS status_id, s.name AS status_name, s.description AS status_description, s.createdAt AS status_createdAt, s.updatedAt AS status_updatedAt, s.deletedAt AS status_deletedAt,
               u.id AS user_id, u.firstName, u.lastName, u.email, u.age, u.lastActiveAt, u.createdAt AS user_createdAt, u.updatedAt AS user_updatedAt, u.deletedAt AS user_deletedAt,
               p.id AS project_id, p.name AS project_name, p.description AS project_description, p.createdAt AS project_createdAt, p.updatedAt AS project_updatedAt, p.deletedAt AS project_deletedAt
        FROM tasks t
        LEFT JOIN statuses s ON s.id = t.statusId
        LEFT JOIN users u ON u.id = t.userId
        LEFT JOIN projects p ON p.id = t.projectId
        WHERE t.deletedAt %v
    `, str)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasksMap := make(map[int]*entities.Task) // Using a map to avoid duplication of tasks
	for rows.Next() {
		task := entities.Task{}

		// Create variables to hold all the values we will scan into
		var status entities.Status
		var user entities.User
		var project entities.Project

		// Scan the result set into the task and associated entities
		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.StatusID, &task.UserID, &task.ProjectID, &task.CreatedAt,
			&task.UpdatedAt, &task.DeletedAt, &task.DeletedBy,
			// Status
			&status.ID, &status.Name, &status.Description, &status.CreatedAt, &status.UpdatedAt, &status.DeletedAt,
			// User
			&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.LastActiveAt, &user.CreatedAt,
			&user.UpdatedAt, &user.DeletedAt,
			// Project
			&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt,
		)

		if err != nil {
			log.Printf("Failed to scan task: %v", err)
			continue
		}

		// Assign the related entities for the task
		task.Status = status
		task.User = user
		task.Project = project

		// Add the task to the map
		tasksMap[task.ID] = &task
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over tasks rows: %v", err)
	}

	// Convert the map to a slice
	tasks := make([]*entities.Task, 0, len(tasksMap))
	for _, task := range tasksMap {
		tasks = append(tasks, task)
	}

	return tasks, nil
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
	err = tx.QueryRow(`
		INSERT INTO tasks (title, description, statusId, userId, projectId) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, title, description, statusId, userId, projectId, createdAt, updatedAt`,
		payload.Title,
		payload.Description,
		payload.StatusID,
		func() interface{} { // Handle optional userId
			if payload.UserID == 0 {
				return nil
			}
			return payload.UserID
		}(),
		payload.ProjectID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.StatusID,
		&task.UserID,
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

func (s *Store) GetTask(id int) (*entities.Task, error) {
	query := fmt.Sprintf(`
        SELECT t.id, t.title, t.description, t.statusId, t.userId, t.projectId, t.createdAt, t.updatedAt, t.deletedAt, t.deletedBy,
			s.id AS status_id, s.name AS status_name, s.description AS status_description, s.createdAt AS status_createdAt, s.updatedAt AS status_updatedAt, s.deletedAt AS status_deletedAt,
			u.id AS user_id, u.firstName, u.lastName, u.email, u.age, u.lastActiveAt, u.createdAt AS user_createdAt, u.updatedAt AS user_updatedAt, u.deletedAt AS user_deletedAt,
			p.id AS project_id, p.name AS project_name, p.description AS project_description, p.createdAt AS project_createdAt, p.updatedAt AS project_updatedAt, p.deletedAt AS project_deletedAt
			FROM tasks t
			LEFT JOIN statuses s ON s.id = t.statusId
			LEFT JOIN users u ON u.id = t.userId
			LEFT JOIN projects p ON p.id = t.projectId
			WHERE t.id = $1 AND t.deletedAt IS NULL
    `)

	row := s.db.QueryRow(query, id)

	task := &entities.Task{}
	var status entities.Status
	var user entities.User
	var project entities.Project

	err := row.Scan(
		&task.ID, &task.Title, &task.Description, &task.StatusID, &task.UserID, &task.ProjectID, &task.CreatedAt,
		&task.UpdatedAt, &task.DeletedAt, &task.DeletedBy,
		// Status
		&status.ID, &status.Name, &status.Description, &status.CreatedAt, &status.UpdatedAt, &status.DeletedAt,
		// User
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.LastActiveAt, &user.CreatedAt,
		&user.UpdatedAt, &user.DeletedAt,
		// Project
		&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve task: %v", err)
	}

	// Assign the related entities for the task
	task.Status = status
	task.User = user
	task.Project = project

	return task, nil
}

func (s *Store) TaskUpdate(payload entities.TaskUpdatePayload) (*entities.Task, error) {
	return nil, nil
}

func (s *Store) TaskDelete(id int) (*entities.Task, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	task := &entities.Task{}
	err = tx.QueryRow("UPDATE tasks SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id, title, description, statusId, userId, projectId, createdAt, updatedAt, deletedAt", id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.StatusID,
		&task.UserID,
		&task.ProjectID,
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
	return nil, nil
}

// TODO: UPDATE TASK (HOW TO UPDATE SUBTASK(ARRAY OF STRINGS))

func scanRowIntoTask(rows *sql.Rows, task *entities.Task) error {
	return rows.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.StatusID,
		&task.UserID,
		&task.ProjectID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
		&task.DeletedBy,
	)
}
