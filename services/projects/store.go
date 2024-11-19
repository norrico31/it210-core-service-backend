package projects

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// SELECT
// 			p.id AS project_id,
// 			p.name AS project_name,
// 			p.description AS project_description,
// 			p.progress AS project_progress,
// 			p.dateStarted AS project_date_started,
// 			p.dateDeadline AS project_date_deadline,
// 			p.createdAt AS project_created_at,
// 			p.updatedAt AS project_updated_at,
// 			p.deletedAt AS project_deleted_at,

// 			u.id AS user_id,
// 			u.firstName AS user_first_name,
// 			u.lastName AS user_last_name,
// 			u.email AS user_email,
// 			u.age AS user_age,

// 			u.lastActiveAt AS user_last_active_at,
// 			u.createdAt AS user_created_at,
// 			u.updatedAt AS user_updated_at,
// 			u.deletedAt AS user_deleted_at,
// 			u.deletedBy AS user_deleted_by,

// 			t.id AS task_id,
// 			t.title AS task_title,
// 			t.description AS task_description,
// 			t.statusId AS task_status_id,
// 			t.userId AS task_user_id,
// 			t.projectId AS task_project_id,
// 			t.createdAt AS task_created_at,
// 			t.updatedAt AS task_updated_at,
// 			t.deletedAt AS task_deleted_at,
// 			t.deletedBy AS task_deleted_by,

// 			s.id AS status_id,
// 			s.name AS status_name,
// 			s.description AS status_description

// 		FROM
// 			projects p
// 		LEFT JOIN
// 			users_projects up ON p.id = up.project_id
// 		LEFT JOIN
// 			users u ON up.user_id = u.id
// 		LEFT JOIN
// 			tasks t ON t.projectId = p.id
// 		LEFT JOIN
// 			statuses s ON s.id = t.statusId
// 		WHERE
// 			p.deletedAt `

// TODO: WHEN DELETE PLEASE INCLUDE THE USER ID TO TRACK WHO'S USER DELETE THE OPERATIONS (deletedBy in tables)
func (s *Store) GetProjects(condition string) ([]*entities.Project, error) {
	query := `
		SELECT
			p.id AS project_id,
			p.name AS project_name,
			p.description AS project_description,
			p.progress AS project_progress,
			p.dateStarted AS project_date_started,
			p.dateDeadline AS project_date_deadline,
			p.createdAt AS project_created_at,
 			p.updatedAt AS project_updated_at,
 			p.deletedAt AS project_deleted_at,
			p.deletedBy project_deleted_by,

			u.id AS user_id,
			u.firstName AS user_first_name,
			u.lastName AS user_last_name,
			u.email AS user_email,
			u.age AS user_age,

			u.lastActiveAt AS user_last_active_at,
			u.createdAt AS user_created_at,
			u.updatedAt AS user_updated_at,
			u.deletedAt AS user_deleted_at,
			u.deletedBy AS user_deleted_by,

			t.id AS task_id,
			t.title AS task_title,
			t.description AS task_description,
			t.statusId AS task_status_id,
			t.userId AS task_user_id,
			t.projectId AS task_project_id,
			t.createdAt AS task_created_at,
			t.updatedAt AS task_updated_at,
			t.deletedAt AS task_deleted_at,
			t.deletedBy AS task_deleted_by,

			s.id status_id,
			s.name status_name,
			s.description status_description

		FROM
			projects p
		LEFT JOIN
			users_projects up ON p.id = up.project_id
		LEFT JOIN
			users u ON up.user_id = u.id
		LEFT JOIN
			tasks t ON t.projectId = p.id
		LEFT JOIN
			statuses s ON s.id = t.statusId
		WHERE
			p.deletedAt ` + condition

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	// Map to hold projects indexed by projectID
	projectsMap := make(map[int]*entities.Project)

	// Loop through query results
	for rows.Next() {
		var projectID int
		var dateStarted, dateDeadline *time.Time
		var project entities.Project
		var user entities.User
		var task entities.Task
		var status entities.Status

		// Use pointers to handle NULL values
		var userID, userAge, userDeletedBy *int
		var userFirstName, userLastName, userEmail *string
		var userLastActiveAt, userCreatedAt, userUpdatedAt, userDeletedAt *time.Time

		// Use pointers for task fields
		var taskID, taskStatusID, taskUserID, taskProjectID, taskDeletedBy *int
		var taskTitle, taskDescription *string
		var taskCreatedAt, taskUpdatedAt, taskDeletedAt *time.Time

		// Scan project, user, and task data
		err = rows.Scan(
			&projectID, &project.Name, &project.Description, &project.Progress, &dateStarted, &dateDeadline, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt, &project.DeletedBy,

			&userID, &userFirstName, &userLastName, &userEmail, &userAge, &userLastActiveAt, &userCreatedAt, &userUpdatedAt, &userDeletedAt, &userDeletedBy,

			&taskID, &taskTitle, &taskDescription, &taskStatusID, &taskUserID, &taskProjectID, &taskCreatedAt, &taskUpdatedAt, &taskDeletedAt, &taskDeletedBy,

			&status.ID, &status.Name, &status.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		// If project is not yet in the map, add it
		if _, exists := projectsMap[projectID]; !exists {
			project.ID = projectID
			project.DateStarted = dateStarted   // Assign dateStarted
			project.DateDeadline = dateDeadline // Assign dateDeadline
			project.Users = []entities.User{}   // Initialize empty users slice
			project.Tasks = []entities.Task{}   // Initialize empty tasks slice
			projectsMap[projectID] = &project
		}

		// Add user data if available
		if userID != nil {
			user.ID = *userID
			if userFirstName != nil {
				user.FirstName = *userFirstName
			}
			if userLastName != nil {
				user.LastName = *userLastName
			}
			if userEmail != nil {
				user.Email = *userEmail
			}
			if userAge != nil {
				user.Age = *userAge
			}

			user.LastActiveAt = userLastActiveAt
			if userCreatedAt != nil {
				user.CreatedAt = *userCreatedAt
			}
			if userUpdatedAt != nil {
				user.UpdatedAt = *userUpdatedAt
			}
			user.DeletedAt = userDeletedAt
			user.DeletedBy = userDeletedBy // This will be `nil` if the database value is `NULL`

			// Add user to the project's user list
			projectsMap[projectID].Users = append(projectsMap[projectID].Users, user)
		}

		// Add task data if available and ensure task is not duplicated for the project
		if taskID != nil {
			// Check if the task is already in the project's tasks
			taskExists := false
			for _, existingTask := range projectsMap[projectID].Tasks {
				if existingTask.ID == *taskID {
					taskExists = true
					break
				}
			}

			// If task doesn't already exist, add it
			if !taskExists {
				task.ID = *taskID
				task.Title = *taskTitle
				task.Description = *taskDescription
				task.StatusID = *taskStatusID
				task.UserID = taskUserID
				task.ProjectID = *taskProjectID
				task.CreatedAt = *taskCreatedAt
				task.UpdatedAt = *taskUpdatedAt
				task.DeletedAt = taskDeletedAt
				task.DeletedBy = taskDeletedBy
				task.Status = status

				// Add task to the project's task list
				projectsMap[projectID].Tasks = append(projectsMap[projectID].Tasks, task)
			}
		}
	}

	// Convert map to slice
	var projects []*entities.Project
	for _, project := range projectsMap {
		projects = append(projects, project)
	}

	// Sort projects by creation date
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].CreatedAt.After(projects[j].CreatedAt)
	})

	return projects, nil
}

func (s *Store) GetProject(id int) (*entities.Project, error) {
	query := `
		SELECT 
			p.id, p.name, p.description, p.createdAt, p.updatedAt, p.deletedAt,

			t.id AS task_id, t.title AS task_title, 
			t.description AS task_description, t.statusId AS task_statusId, 
			t.userId AS task_userId, t.createdAt AS task_createdAt, 
			t.updatedAt AS task_updatedAt, t.deletedAt AS task_deletedAt, t.deletedBy AS task_deletedBy,

			s.id AS status_id, s.name AS status_name, s.description AS status_description, 
			s.createdAt AS status_createdAt, s.updatedAt AS status_updatedAt, s.deletedAt AS status_deletedAt,

			u.id AS user_id, u.firstName AS user_firstName, u.lastName AS user_lastName, 
			u.email AS user_email, u.age AS user_age, u.lastActiveAt, 
			u.createdAt AS user_createdAt, u.updatedAt AS user_updatedAt, 
			u.deletedAt AS user_deletedAt, u.deletedBy AS user_deletedBy,

			up_user.id AS project_user_id, up_user.firstName AS project_user_first_name, 
			up_user.lastName AS project_user_last_name, up_user.email AS project_user_email, 
			up_user.age AS project_user_age, up_user.lastActiveAt AS project_user_last_active_at,
			up_user.createdAt AS project_user_created_at, up_user.updatedAt AS project_user_updated_at, 
			up_user.deletedAt AS project_user_deleted_at, up_user.deletedBy AS project_user_deleted_by

		FROM 
			projects p
		LEFT JOIN 
			tasks t ON t.projectId = p.id
		LEFT JOIN 
			statuses s ON t.statusId = s.id
		LEFT JOIN 
			users u ON t.userId = u.id
		LEFT JOIN 
			users_projects up ON up.project_id = p.id
		LEFT JOIN 
			users up_user ON up.user_id = up_user.id
		WHERE 
			p.id = $1 AND p.deletedAt IS NULL
	`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	project := entities.Project{}
	taskMap := make(map[int]entities.Task)
	userMap := make(map[int]entities.User)

	for rows.Next() {
		var (
			taskID, statusID, userID, projectUserID            int
			taskDeletedBy, userDeletedBy, projectUserDeletedBy *int
			task                                               entities.Task
			status                                             entities.Status
			user                                               entities.User
			projectUser                                        entities.User
			taskDeletedAt, userDeletedAt, projectUserDeletedAt *time.Time
		)

		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt,
			&project.Progress,

			&taskID, &task.Title, &task.Description, &task.StatusID,
			&task.UserID, &task.CreatedAt, &task.UpdatedAt, &task.DeletedAt, &taskDeletedBy,

			&statusID, &status.Name, &status.Description, &status.CreatedAt,
			&status.UpdatedAt, &status.DeletedAt,

			&userID, &user.FirstName, &user.LastName, &user.Email, &user.Age,
			&user.LastActiveAt, &user.CreatedAt, &user.UpdatedAt, &userDeletedAt, &userDeletedBy,

			&projectUserID, &projectUser.FirstName, &projectUser.LastName,
			&projectUser.Email, &projectUser.Age, &projectUser.LastActiveAt,
			&projectUser.CreatedAt, &projectUser.UpdatedAt, &projectUserDeletedAt, &projectUserDeletedBy,
		)
		if err != nil {
			return nil, err
		}

		// Add task
		if taskID != 0 {
			if _, exists := taskMap[taskID]; !exists {
				task.ID = taskID
				if statusID != 0 {
					status.ID = statusID
					task.Status = status
				}
				if userID != 0 {
					user.ID = userID
					user.DeletedAt = userDeletedAt
					user.DeletedBy = userDeletedBy
					task.User = user
				}
				task.DeletedAt = taskDeletedAt
				task.DeletedBy = taskDeletedBy
				taskMap[taskID] = task
			}
		}

		// Add project user
		if projectUserID != 0 {
			if _, exists := userMap[projectUserID]; !exists {
				projectUser.ID = projectUserID
				projectUser.DeletedAt = projectUserDeletedAt
				projectUser.DeletedBy = projectUserDeletedBy
				userMap[projectUserID] = projectUser
			}
		}
	}

	// Collect tasks
	for _, task := range taskMap {
		project.Tasks = append(project.Tasks, task)
	}

	// Collect users
	for _, user := range userMap {
		project.Users = append(project.Users, user)
	}

	if project.ID == 0 {
		return nil, fmt.Errorf("project not found")
	}

	return &project, nil
}

func (s *Store) ProjectCreate(payload entities.ProjectCreatePayload) (*entities.Project, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	proj := entities.Project{}
	err = tx.QueryRow("INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id, name, description, createdAt, updatedAt", payload.Name, payload.Description).Scan(
		&proj.ID,
		&proj.Name,
		&proj.Description,
		&proj.CreatedAt,
		&proj.UpdatedAt,
	)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rollbackErr)
		}
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &proj, err
}

func (s *Store) ProjectUpdate(payload entities.ProjectUpdatePayload) (*entities.Project, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	proj := &entities.Project{}
	err = tx.QueryRow("UPDATE projects SET name = $1, description = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, name, description, createdAt, updatedAt", payload.Name, payload.Description, payload.ID).Scan(
		&proj.ID,
		&proj.Name,
		&proj.Description,
		&proj.CreatedAt,
		&proj.UpdatedAt,
	)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rollbackErr)
		}
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return proj, err
}

func (s *Store) ProjectDelete(id int) (*entities.Project, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	proj := entities.Project{}
	err = tx.QueryRow("UPDATE projects SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id, name, description, createdAt, updatedAt, deletedAt", id).Scan(
		&proj.ID,
		&proj.Name,
		&proj.Description,
		&proj.CreatedAt,
		&proj.UpdatedAt,
		&proj.DeletedAt,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error deleting: %v rollback error: %v", err, rollbackErr)
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &proj, err
}

func (s *Store) ProjectRestore(id int) (*entities.Project, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	fmt.Println("hala?")
	proj := entities.Project{}
	err = tx.QueryRow("UPDATE projects SET deletedAt = NULL WHERE id = $1 RETURNING id, name, description, createdAt, updatedAt, deletedAt", id).Scan(
		&proj.ID,
		&proj.Name,
		&proj.Description,
		&proj.CreatedAt,
		&proj.UpdatedAt,
		&proj.DeletedAt,
	)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error restoring: %v rollback error: %v", err, rollbackErr)
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return &proj, nil
	}

	return &proj, nil
}

func scanRowIntoProject(rows *sql.Rows, proj *entities.Project) error {
	return rows.Scan(
		&proj.ID,
		&proj.Name,
		&proj.Description,
		&proj.CreatedAt,
		&proj.UpdatedAt,
		&proj.DeletedAt,
	)
}
