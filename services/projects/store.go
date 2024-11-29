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
			u.roleId user_role_id,
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
			s.description status_description,

			ut.id AS user_task_id,
			ut.firstName AS user_task_first_name,
			ut.lastName AS user_task_last_name,
			ut.email AS user_task_email,
			ut.age AS user_task_age,
			ut.lastActiveAt AS user_task_last_active_at,
			ut.createdAt AS user_task_created_at,
			ut.updatedAt AS user_task_updated_at,
			ut.deletedAt AS user_task_deleted_at,
			ut.deletedBy AS user_task_deleted_by,

			r.id role_id,
			r.name role_name,
			r.description role_description,
			r.createdAt role_created_at,
			r.updatedAt role_updated_at,
			r.deletedAt role_deleted_at

		FROM
			projects p
		LEFT JOIN
			users_projects up ON up.deletedAt IS NULL AND p.id = up.project_id
		LEFT JOIN
			users u ON up.deletedAt IS NULL AND up.user_id = u.id
		LEFT JOIN
			tasks t ON t.deletedAt IS NULL AND t.projectId = p.id
		LEFT JOIN
			statuses s ON s.deletedAt IS NULL AND s.id = t.statusId
		LEFT JOIN
			users ut ON ut.deletedAt IS NULL AND ut.id = t.userId
		LEFT JOIN
			roles r ON r.deletedAt IS NULL AND r.id = u.roleId
		WHERE
			p.deletedAt ` + condition

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	projectsMap := make(map[int]*entities.Project)

	// Loop through query results
	for rows.Next() {
		var projectID int
		// var dateStarted, dateDeadline *time.Time
		var project entities.Project

		var user entities.User
		var task entities.Task
		var status entities.Status
		var statusID *int
		var statusName, statusDescription *string

		// Use pointers to handle NULL values
		var userID, userAge, userRoleId, userDeletedBy *int
		var userFirstName, userLastName, userEmail *string
		var userLastActiveAt, userCreatedAt, userUpdatedAt, userDeletedAt *time.Time

		// Use pointers for task fields
		var taskID, taskStatusID, taskUserID, taskProjectID, taskDeletedBy *int
		var taskTitle, taskDescription *string
		var taskCreatedAt, taskUpdatedAt, taskDeletedAt *time.Time

		var userIDTask, taskUserAge, taskUserDeletedBy *int
		var taskUserFirstName, taskUserLastName, taskUserEmail *string
		var taskUserLastActiveAt, taskUserCreatedAt, taskUserUpdatedAt, taskUserDeletedAt *time.Time

		var roleId *int
		var roleName, roleDescription *string
		var roleCreatedAt, roleUpdatedAt, roleDeletedAt *time.Time

		taskUser := entities.User{}

		// Scan project, user, and task data
		err = rows.Scan(
			&projectID, &project.Name, &project.Description, &project.Progress, &project.DateStarted, &project.DateDeadline, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt, &project.DeletedBy,

			&userID, &userFirstName, &userLastName, &userEmail, &userAge, &userRoleId, &userLastActiveAt, &userCreatedAt, &userUpdatedAt, &userDeletedAt, &userDeletedBy,

			&taskID, &taskTitle, &taskDescription, &taskStatusID, &taskUserID, &taskProjectID, &taskCreatedAt, &taskUpdatedAt, &taskDeletedAt, &taskDeletedBy,

			&statusID, &statusName, &statusDescription,

			&userIDTask, &taskUserFirstName, &taskUserLastName, &taskUserEmail, &taskUserAge, &taskUserLastActiveAt, &taskUserCreatedAt, &taskUserUpdatedAt, &taskUserDeletedAt, &taskUserDeletedBy,

			&roleId, &roleName, &roleDescription, &roleCreatedAt, &roleUpdatedAt, &roleDeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		// If project is not yet in the map, add it
		if _, exists := projectsMap[projectID]; !exists {
			project.ID = projectID
			project.Users = []entities.User{} // Initialize empty users slice
			project.Tasks = []entities.Task{} // Initialize empty tasks slice
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
			if userRoleId != nil {
				user.RoleId = roleId

				role := entities.Role{
					ID:          *roleId,
					Name:        *roleName,
					Description: *roleDescription,
					CreatedAt:   *roleCreatedAt,
					UpdatedAt:   *roleUpdatedAt,
					DeletedAt:   roleDeletedAt,
				}
				user.Role = role
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

		if statusID != nil {
			status.ID = *statusID
		}

		if statusName != nil {
			status.Name = *statusName
		}

		if statusDescription != nil {
			status.Description = *statusDescription
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

			// debug this
			if taskUserID != nil {
				taskUser.ID = *userIDTask
				if userFirstName != nil {
					taskUser.FirstName = *taskUserFirstName
				}
				if userLastName != nil {
					taskUser.LastName = *taskUserLastName
				}
				if userEmail != nil {
					taskUser.Email = *taskUserEmail
				}
				if userAge != nil {
					taskUser.Age = *taskUserAge
				}

				taskUser.LastActiveAt = taskUserLastActiveAt
				if taskUserCreatedAt != nil {
					taskUser.CreatedAt = *taskUserCreatedAt
				}
				if taskUserUpdatedAt != nil {
					taskUser.UpdatedAt = *taskUserUpdatedAt
				}
				taskUser.DeletedAt = taskUserDeletedAt
				taskUser.DeletedBy = taskUserDeletedBy
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
				task.User = taskUser

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

// TODO: REFACTOR SEPARATE THE ARRAYS INTO OTHER QUERY
func (s *Store) GetProject(id int) (*entities.Project, error) {
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
			p.deletedBy AS project_deleted_by,

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
			s.description status_description,
			ut.id AS user_task_id,
			ut.firstName AS user_task_first_name,
			ut.lastName AS user_task_last_name,
			ut.email AS user_task_email,
			ut.age AS user_task_age,
			ut.lastActiveAt AS user_task_last_active_at,
			ut.createdAt AS user_task_created_at,
			ut.updatedAt AS user_task_updated_at,
			ut.deletedAt AS user_task_deleted_at,
			ut.deletedBy AS user_task_deleted_by
			
		FROM 
			projects p
		LEFT JOIN
			users_projects up ON p.id = up.project_id
		LEFT JOIN
			users u ON up.user_id = u.id
		LEFT JOIN
			tasks t ON t.deletedAt IS NULL AND t.projectId = p.id
		LEFT JOIN
			statuses s ON s.id = t.statusId
		LEFT JOIN
			users ut ON ut.id = t.userId
		WHERE 
			p.id = $1 AND p.deletedAt IS NULL
	`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	project := entities.Project{}

	for rows.Next() {
		var dateStarted, dateDeadline *time.Time
		user := entities.User{}
		var userID, userAge, userDeletedBy *int
		var userFirstName, userLastName, userEmail *string
		var userLastActiveAt, userCreatedAt, userUpdatedAt, userDeletedAt *time.Time

		task := entities.Task{}
		var taskID, taskStatusID, taskUserID, taskProjectID, taskDeletedBy *int
		var taskTitle, taskDescription *string
		var taskCreatedAt, taskUpdatedAt, taskDeletedAt *time.Time

		taskUser := entities.User{}
		var userIDTask, taskUserAge, taskUserDeletedBy *int
		var taskUserFirstName, taskUserLastName, taskUserEmail *string
		var taskUserLastActiveAt, taskUserCreatedAt, taskUserUpdatedAt, taskUserDeletedAt *time.Time

		status := entities.Status{}
		var statusID *int
		var statusName, statusDescription *string

		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Progress, &dateStarted, &dateDeadline, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt, &project.DeletedBy,
			&userID, &userFirstName, &userLastName, &userEmail, &userAge, &userLastActiveAt, &userCreatedAt, &userUpdatedAt, &userDeletedAt, &userDeletedBy,
			&taskID, &taskTitle, &taskDescription, &taskStatusID, &taskUserID, &taskProjectID, &taskCreatedAt, &taskUpdatedAt, &taskDeletedAt, &taskDeletedBy,
			&statusID, &statusName, &statusDescription,
			&userIDTask, &taskUserFirstName, &taskUserLastName, &taskUserEmail, &taskUserAge, &taskUserLastActiveAt, &taskUserCreatedAt, &taskUserUpdatedAt, &taskUserDeletedAt, &taskUserDeletedBy,
		)
		if err != nil {
			return nil, err
		}

		if dateStarted != nil {
			project.DateStarted = dateStarted
		}
		if dateDeadline != nil {
			project.DateDeadline = dateDeadline
		}

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
			project.Users = append(project.Users, user)
		}

		if statusID != nil {
			status.ID = *statusID
		}

		if statusName != nil {
			status.Name = *statusName
		}

		if statusDescription != nil {
			status.Description = *statusDescription
		}

		if taskID != nil {
			// Check if the task is already in the project's tasks
			taskExists := false
			for _, existingTask := range project.Tasks {
				if existingTask.ID == *taskID {
					taskExists = true
					break
				}
			}

			// debug this
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
				if taskUserID != nil {
					taskUser.ID = *userIDTask
					if userFirstName != nil {
						taskUser.FirstName = *taskUserFirstName
					}
					if userLastName != nil {
						taskUser.LastName = *taskUserLastName
					}
					if userEmail != nil {
						taskUser.Email = *taskUserEmail
					}
					if userAge != nil {
						taskUser.Age = *taskUserAge
					}

					taskUser.LastActiveAt = taskUserLastActiveAt
					if taskUserCreatedAt != nil {
						taskUser.CreatedAt = *taskUserCreatedAt
					}
					if taskUserUpdatedAt != nil {
						taskUser.UpdatedAt = *taskUserUpdatedAt
					}
					taskUser.DeletedAt = taskUserDeletedAt
					taskUser.DeletedBy = taskUserDeletedBy
				}

				task.User = taskUser

				// Add task to the project's task list
				project.Tasks = append(project.Tasks, task)
			}
		}

	}

	if project.ID == 0 {
		return nil, fmt.Errorf("project not found")
	}

	return &project, nil
}

func (s *Store) ProjectCreate(payload entities.ProjectCreatePayload) (map[string]interface{}, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}
	progress := 0.0

	if payload.Progress != nil {
		progress = *payload.Progress
	}

	var started, deadline *time.Time

	if payload.DateStarted != "" {
		dateStarted, err := normalizeDate(payload.DateStarted)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for DateStarted: %v", err)
		}
		started = &dateStarted
	}
	if payload.DateDeadline != "" {
		dateDeadline, err := normalizeDate(payload.DateDeadline)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for DateDeadline: %v", err)
		}
		deadline = &dateDeadline
	}

	proj := entities.Project{}
	err = tx.QueryRow("INSERT INTO projects (name, description, progress, dateStarted, dateDeadline) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, description, progress, dateStarted, dateDeadline, createdAt, updatedAt", payload.Name, payload.Description, progress, started, deadline).Scan(
		&proj.ID,
		&proj.Name,
		&proj.Description,
		&proj.Progress,
		&proj.DateStarted,
		&proj.DateDeadline,
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

	return buildProjectResponse(proj), err
}

func (s *Store) ProjectUpdate(payload entities.ProjectUpdatePayload) (map[string]interface{}, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	// Fetch the existing project
	query := `
		SELECT id, name, description, progress, dateStarted, dateDeadline, createdAt, updatedAt
		FROM projects
		WHERE id = $1`
	rows, err := s.db.Query(query, payload.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching project: %v", err)
	}
	defer rows.Close()

	existProj := entities.Project{}
	if rows.Next() {
		err := rows.Scan(
			&existProj.ID, &existProj.Name, &existProj.Description, &existProj.Progress,
			&existProj.DateStarted, &existProj.DateDeadline, &existProj.CreatedAt, &existProj.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	// Update project fields based on the payload
	if payload.Name != "" {
		existProj.Name = payload.Name
	}
	if payload.Description != "" {
		existProj.Description = payload.Description
	}
	if payload.Progress != nil {
		existProj.Progress = payload.Progress
	}
	if payload.DateStarted != "" {
		dateStarted, err := normalizeDate(payload.DateStarted)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for DateStarted: %v", err)
		}
		existProj.DateStarted = &dateStarted
	}
	if payload.DateDeadline != "" {
		dateDeadline, err := normalizeDate(payload.DateDeadline)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for DateDeadline: %v", err)
		}
		existProj.DateDeadline = &dateDeadline
	}

	updateQuery := `
		UPDATE projects
		SET name = $1, description = $2, progress = $3, dateStarted = $4, dateDeadline = $5, updatedAt = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, name, description, progress, dateStarted, dateDeadline, createdAt, updatedAt`
	proj := entities.Project{}
	err = tx.QueryRow(updateQuery,
		existProj.Name,
		existProj.Description,
		existProj.Progress,
		existProj.DateStarted,
		existProj.DateDeadline,
		existProj.ID,
	).Scan(
		&proj.ID, &proj.Name, &proj.Description, &proj.Progress,
		&proj.DateStarted, &proj.DateDeadline, &proj.CreatedAt, &proj.UpdatedAt,
	)

	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("update error: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Build and return the response
	return buildProjectResponse(proj), nil
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

func normalizeDate(input string) (time.Time, error) {
	formats := []string{"1/2/2006", "01/02/2006"} // Support both single and double-digit formats
	var parsedDate time.Time
	var err error
	for _, format := range formats {
		parsedDate, err = time.Parse(format, input)
		if err == nil {
			return parsedDate, nil // Return the first successfully parsed date
		}
	}
	return time.Time{}, fmt.Errorf("could not parse date: %v", input)
}

func buildProjectResponse(proj entities.Project) map[string]interface{} {
	formatDate := func(date *time.Time) interface{} {
		if date != nil {
			return date.Format("01/02/2006") // MM/DD/YYYY
		}
		return nil
	}

	return map[string]interface{}{
		"id":           proj.ID,
		"name":         proj.Name,
		"description":  proj.Description,
		"progress":     proj.Progress,
		"dateStarted":  formatDate(proj.DateStarted),
		"dateDeadline": formatDate(proj.DateDeadline),
		"createdAt":    proj.CreatedAt,
		"updatedAt":    proj.UpdatedAt,
	}
}
