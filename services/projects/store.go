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
			p.url AS project_url,
			p.progress AS project_progress,
			p.statusId AS project_status_id,
			p.dateStarted AS project_date_started,
			p.dateDeadline AS project_date_deadline,
			p.createdAt AS project_created_at,
			p.updatedAt AS project_updated_at,
			p.deletedAt AS project_deleted_at,
			p.deletedBy AS project_deleted_by,

			stat.id AS status_id,
			stat.name AS status_name,
			stat.description AS status_description,

			seg.id segment_id,
			seg.name segment_name,
			seg.description segment_description,

			u.id AS user_id,
			u.firstName AS user_first_name,
			u.lastName AS user_last_name,
			u.email AS user_email,
			u.age AS user_age,
			u.roleId AS user_role_id,
			u.lastActiveAt AS user_last_active_at,
			u.createdAt AS user_created_at,
			u.updatedAt AS user_updated_at,
			u.deletedAt AS user_deleted_at,
			u.deletedBy AS user_deleted_by,

			t.id AS task_id,
			t.name AS task_name,
			t.description AS task_description,
			t.userId AS task_user_id,
			t.priorityId AS task_priority_id,
			t.createdAt AS task_created_at,
			t.updatedAt AS task_updated_at,
			t.deletedAt AS task_deleted_at,
			t.deletedBy AS task_deleted_by

		FROM 
			projects p
		LEFT JOIN
			users_projects up ON up.deletedAt IS NULL AND p.id = up.project_id
		LEFT JOIN
			segments_projects sp ON sp.projectId = p.id
		LEFT JOIN
			segments seg ON seg.id = sp.segmentId
		LEFT JOIN
			users u ON up.deletedAt IS NULL AND up.user_id = u.id
		LEFT JOIN
			statuses stat ON stat.id = p.statusId
		LEFT JOIN
			project_tasks t ON t.deletedAt IS NULL AND t.projectId = p.id
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projectsMap := make(map[int]*entities.Project)
	userMap := make(map[int]bool)

	for rows.Next() {
		var project entities.Project
		var dateStarted, dateDeadline *time.Time
		var projectId, projectStatusId *int

		var statusID *int
		var statusName, statusDescription *string

		user := entities.User{}
		var userFirstName, userLastName, userEmail, segmentName, segmentDescription *string
		var userID, userAge, userDeletedBy, userRoleId, segmentId *int
		var userLastActiveAt, userCreatedAt, userUpdatedAt, userDeletedAt *time.Time

		task := entities.TasksProject{}
		var taskID, taskUserID, taskPriorityID *int
		var taskName, taskDescription *string
		var taskCreatedAt, taskUpdatedAt, taskDeletedAt *time.Time
		var taskDeletedBy *int

		err := rows.Scan(
			&projectId, &project.Name, &project.Description, &project.Url, &project.Progress, &projectStatusId, &dateStarted, &dateDeadline, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt, &project.DeletedBy,
			&statusID, &statusName, &statusDescription,
			&segmentId, &segmentName, &segmentDescription,
			&userID, &userFirstName, &userLastName, &userEmail, &userAge, &userRoleId, &userLastActiveAt, &userCreatedAt, &userUpdatedAt, &userDeletedAt, &userDeletedBy,
			&taskID, &taskName, &taskDescription, &taskUserID, &taskPriorityID, &taskCreatedAt, &taskUpdatedAt, &taskDeletedAt, &taskDeletedBy,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := projectsMap[*projectId]; !exists {
			project.ID = *projectId
			project.Users = []entities.User{}
			project.Tasks = []entities.TasksProject{}
			projectsMap[*projectId] = &project
		}

		if userID != nil && !userMap[*userID] {
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
				user.RoleId = userRoleId
			}

			user.LastActiveAt = userLastActiveAt
			if userCreatedAt != nil {
				user.CreatedAt = *userCreatedAt
			}
			if userUpdatedAt != nil {
				user.UpdatedAt = *userUpdatedAt
			}
			user.DeletedAt = userDeletedAt
			user.DeletedBy = userDeletedBy

			userMap[*userID] = true
			projectsMap[*projectId].Users = append(projectsMap[*projectId].Users, user)
		}

		if taskID != nil {
			task.ID = *taskID
			if taskName != nil {
				task.Name = *taskName
			}
			if taskDescription != nil {
				task.Description = *taskDescription
			}
			task.UserID = taskUserID
			task.PriorityID = *taskPriorityID
			task.CreatedAt = *taskCreatedAt
			task.UpdatedAt = *taskUpdatedAt
			task.DeletedAt = taskDeletedAt
			task.DeletedBy = taskDeletedBy

			projectsMap[*projectId].Tasks = append(projectsMap[*projectId].Tasks, task)
		}

		if projectStatusId != nil {
			project.StatusID = *projectStatusId
			status := entities.Status{
				ID:          *projectStatusId,
				Name:        *statusName,
				Description: *statusDescription,
			}
			project.Status = status
		}

		if segmentId != nil {
			project.SegmentID = *segmentId
			segment := entities.Segment{
				ID:          *segmentId,
				Name:        *segmentName,
				Description: *segmentDescription,
			}
			project.Segment = segment
		}

		if dateStarted != nil {
			project.DateStarted = dateStarted
		}
		if dateDeadline != nil {
			project.DateDeadline = dateDeadline
		}
	}

	var projects []*entities.Project
	for _, project := range projectsMap {
		projects = append(projects, project)
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].CreatedAt.After(projects[j].CreatedAt)
	})

	return projects, nil
}

func (s *Store) GetProject(id int) (*entities.Project, error) {
	// Query to retrieve project and related data
	projectQuery := `
		SELECT 
			p.id AS project_id,
			p.name AS project_name,
			p.description AS project_description,
			p.url AS project_url,
			p.progress AS project_progress,
			p.statusId AS project_status_id,
			p.dateStarted AS project_date_started,
			p.dateDeadline AS project_date_deadline,
			p.createdAt AS project_created_at,
			p.updatedAt AS project_updated_at,
			p.deletedAt AS project_deleted_at,
			p.deletedBy AS project_deleted_by,

			stat.id AS status_id,
			stat.name AS status_name,
			stat.description AS status_description,

			seg.id AS segment_id,
			seg.name AS segment_name,
			seg.description AS segment_description,

			u.id AS user_id,
			u.firstName AS user_first_name,
			u.lastName AS user_last_name,
			u.email AS user_email,
			u.age AS user_age,
			u.roleId AS user_role_id,
			u.lastActiveAt AS user_last_active_at,
			u.createdAt AS user_created_at,
			u.updatedAt AS user_updated_at,
			u.deletedAt AS user_deleted_at,
			u.deletedBy AS user_deleted_by

		FROM 
			projects p
		LEFT JOIN
			users_projects up ON p.id = up.project_id
		LEFT JOIN
			segments_projects sp ON sp.projectId = p.id
		LEFT JOIN
			segments seg ON seg.id = sp.segmentId
		LEFT JOIN
			users u ON up.user_id = u.id
		LEFT JOIN
			statuses stat ON stat.id = p.statusId
		WHERE 
			p.id = $1 AND p.deletedAt IS NULL
	`

	// Query to retrieve tasks associated with the project, including user and priority details
	tasksQuery := `
		SELECT 
			t.id, 
			t.name, 
			t.description, 
			t.userId, 
			t.priorityId, 
			t.projectId, 
			t.createdAt, 
			t.updatedAt, 
			t.deletedAt, 
			t.deletedBy,
			u.id AS user_id,
			u.firstName AS user_first_name,
			u.lastName AS user_last_name,
			u.email AS user_email,
			p.id AS priority_id,
			p.name AS priority_name,
			p.description AS priority_description
		FROM 
			project_tasks t
		LEFT JOIN 
			users u ON t.userId = u.id
		LEFT JOIN 
			priorities p ON t.priorityId = p.id
		WHERE 
			t.projectId = $1 AND t.deletedAt IS NULL
	`

	// Fetch project details
	rows, err := s.db.Query(projectQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	project := entities.Project{}
	userMap := make(map[int]bool)

	for rows.Next() {
		var dateStarted, dateDeadline *time.Time
		user := entities.User{}
		var userID, userAge, userDeletedBy, projectStatusId, userRoleId *int
		var userFirstName, userLastName, userEmail *string
		var userLastActiveAt, userCreatedAt, userUpdatedAt, userDeletedAt *time.Time

		var statusID, segmentId *int
		var statusName, statusDescription *string
		var segmentName, segmentDescription *string

		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Url, &project.Progress, &projectStatusId, &dateStarted, &dateDeadline, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt, &project.DeletedBy,
			&statusID, &statusName, &statusDescription, &segmentId, &segmentName, &segmentDescription,
			&userID, &userFirstName, &userLastName, &userEmail, &userAge, &userRoleId, &userLastActiveAt, &userCreatedAt, &userUpdatedAt, &userDeletedAt, &userDeletedBy,
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

		if statusID != nil {
			project.StatusID = *statusID
			project.Status = entities.Status{
				ID:          *statusID,
				Name:        *statusName,
				Description: *statusDescription,
			}
		}

		if segmentId != nil {
			project.SegmentID = *segmentId
			project.Segment = entities.Segment{
				ID:          *segmentId,
				Name:        *segmentName,
				Description: *segmentDescription,
			}
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
			user.DeletedBy = userDeletedBy

			if _, exists := userMap[user.ID]; !exists {
				userMap[user.ID] = true
				project.Users = append(project.Users, user)
			}
		}
	}

	if project.ID == 0 {
		return nil, fmt.Errorf("project not found")
	}

	// Fetch project tasks with user and priority details
	taskRows, err := s.db.Query(tasksQuery, id)
	if err != nil {
		return nil, err
	}
	defer taskRows.Close()

	taskMap := make(map[int]bool)
	for taskRows.Next() {
		task := entities.TasksProject{}
		user := entities.User{}
		priority := entities.Priority{}
		var taskDeletedAt *time.Time
		var taskDeletedBy *int

		err := taskRows.Scan(
			&task.ID, &task.Name, &task.Description, &task.UserID, &task.PriorityID, &task.ProjectID,
			&task.CreatedAt, &task.UpdatedAt, &taskDeletedAt, &taskDeletedBy,
			&user.ID, &user.FirstName, &user.LastName, &user.Email,
			&priority.ID, &priority.Name, &priority.Description,
		)
		if err != nil {
			return nil, err
		}

		task.User = user
		task.Priority = priority
		task.DeletedAt = taskDeletedAt
		task.DeletedBy = taskDeletedBy

		if !taskMap[task.ID] {
			taskMap[task.ID] = true
			project.Tasks = append(project.Tasks, task)
		}
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
	err = tx.QueryRow(`
		INSERT INTO projects (name, description, progress, url, statusId, dateStarted, dateDeadline, createdAt, updatedAt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id, name, description, progress, url, statusId, dateStarted, dateDeadline, createdAt, updatedAt`,
		payload.Name,
		payload.Description,
		progress,
		payload.Url,
		payload.StatusID,
		started,
		deadline,
		time.Now(),
		time.Now(),
	).Scan(
		&proj.ID,
		&proj.Name,
		&proj.Description,
		&proj.Progress,
		&proj.Url,
		&proj.StatusID,
		&proj.DateStarted,
		&proj.DateDeadline,
		&proj.CreatedAt,
		&proj.UpdatedAt,
	)

	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create project: %v", err)
	}

	if proj.ID == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create project, project ID is invalid")
	}

	if len(*payload.UserIDs) > 0 {
		for _, userID := range *payload.UserIDs {
			_, err = tx.Exec(`
				INSERT INTO users_projects (user_id, project_id)
				VALUES ($1, $2)
			`, userID, proj.ID)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to associate user with project %d: %v", proj.ID, err)
			}
		}
	}

	if payload.SegmentID != nil {
		_, err = tx.Exec(`
		INSERT INTO segments_projects (segmentId, projectId, deletedAt, deletedBy)
		VALUES ($1, $2, NULL, NULL)
	`, payload.SegmentID, proj.ID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to associate segment with project: %v", err)
		}
	}

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rollbackErr)
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return buildProjectResponse(proj), nil
}

func (s *Store) ProjectUpdate(projId int, payload entities.ProjectUpdatePayload, userIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	var dateStarted, dateDeadline *string
	if payload.DateStarted != "" {
		dateStarted = &payload.DateStarted
	}
	if payload.DateDeadline != "" {
		dateDeadline = &payload.DateDeadline
	}

	updateQuery := `
		UPDATE projects
		SET name = $1, description = $2, progress = $3, url = $4, dateStarted = $5, dateDeadline = $6, statusId = $7, updatedAt = $8
		WHERE id = $9
		RETURNING id, name, description, progress, url, dateStarted, dateDeadline, statusId, createdAt, updatedAt`
	_, err = tx.Exec(updateQuery,
		payload.Name,
		payload.Description,
		payload.Progress,
		payload.Url,
		dateStarted,
		dateDeadline,
		payload.StatusID,
		time.Now(),
		projId,
	)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("update error: %v", err)
	}

	deleteQuery := `DELETE FROM segments_projects WHERE projectId = $1 RETURNING projectId, segmentId`
	rows, err := tx.Query(deleteQuery, projId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete old segment associations: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var projectId, segmentId int
		if err := rows.Scan(&projectId, &segmentId); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to read deleted segment association: %v", err)
		}
		fmt.Printf("Deleted segment association: projectId=%d, segmentId=%d\n", projectId, segmentId)
	}

	_, err = tx.Exec(`
			INSERT INTO segments_projects (segmentId, projectId, deletedAt, deletedBy)
			VALUES ($1, $2, NULL, NULL)
		`, payload.SegmentID, projId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to associate segment with project: %v", err)
	}

	_, err = tx.Exec(`
		DELETE FROM users_projects WHERE project_id = $1
	`, projId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete old user-project associations: %v", err)
	}

	for _, userID := range userIDs {
		_, err = tx.Exec(`
			INSERT INTO users_projects (user_id, project_id)
			VALUES ($1, $2)
		`, userID, projId)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to associate user with project %d: %v", projId, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *Store) ProjectDelete(id int) (*entities.Project, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	println("projectId: ", id)
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
		"statusId":     proj.StatusID,
		"url":          proj.Url,
		"dateStarted":  formatDate(proj.DateStarted),
		"dateDeadline": formatDate(proj.DateDeadline),
		"createdAt":    proj.CreatedAt,
		"updatedAt":    proj.UpdatedAt,
	}
}
