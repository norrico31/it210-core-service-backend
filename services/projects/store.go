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
			p.segmentId AS project_segment_id,
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
			u.deletedBy AS user_deleted_by

		FROM 
			projects p
		LEFT JOIN
			users_projects up ON up.deletedAt IS NULL AND p.id = up.project_id
		LEFT JOIN
			users u ON up.deletedAt IS NULL AND up.user_id = u.id
		LEFT JOIN
			statuses stat ON stat.id = p.statusId
		LEFT JOIN
			segments seg ON seg.id = p.segmentId
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
		var projectId, projectSegmentId, projectStatusId *int

		var statusID, segmentId *int
		var statusName, statusDescription, segmentName, segmentDescription *string

		user := entities.User{}
		var userFirstName, userLastName, userEmail *string
		var userID, userAge, userDeletedBy, userRoleId *int
		var userLastActiveAt, userCreatedAt, userUpdatedAt, userDeletedAt *time.Time

		err := rows.Scan(
			&projectId, &project.Name, &project.Description, &project.Url, &project.Progress, &projectSegmentId, &projectStatusId, &dateStarted, &dateDeadline, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt, &project.DeletedBy,
			&statusID, &statusName, &statusDescription,
			&segmentId, &segmentName, &segmentDescription,
			&userID, &userFirstName, &userLastName, &userEmail, &userAge, &userRoleId, &userLastActiveAt, &userCreatedAt, &userUpdatedAt, &userDeletedAt, &userDeletedBy,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := projectsMap[*projectId]; !exists {
			project.ID = *projectId
			project.Users = []entities.User{}
			projectsMap[*projectId] = &project
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

			if _, exists := userMap[user.ID]; !exists {
				userMap[user.ID] = true
				project.Users = append(project.Users, user)
			}
			projectsMap[*projectId].Users = append(projectsMap[*projectId].Users, user)
		}

		if projectStatusId != nil {
			project.StatusID = *projectStatusId
			status := entities.Status{
				ID:          *projectStatusId,
				Name:        *statusName,
				Description: *statusDescription,
			}
			project.StatusID = *projectStatusId
			project.Status = status
		}

		if projectSegmentId != nil {
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
	query := `
		SELECT 
			p.id AS project_id,
			p.name AS project_name,
			p.description AS project_description,
			p.url AS project_url,
			p.progress AS project_progress,
			p.segmentId project_segment_id,
			p.statusId project_status_id,
			p.dateStarted AS project_date_started,
			p.dateDeadline AS project_date_deadline,
			p.createdAt AS project_created_at,
 			p.updatedAt AS project_updated_at,
 			p.deletedAt AS project_deleted_at,
			p.deletedBy AS project_deleted_by,

			stat.id status_id,
			stat.name status_name,
			stat.description status_description,

			seg.id segment_id,
			seg.name segment_name,
			seg.description segment_description,

			u.id AS user_id,
			u.firstName AS user_first_name,
			u.lastName AS user_last_name,
			u.email AS user_email,
			u.age AS user_age,
			u.roleId as user_role_id,
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
			users u ON up.user_id = u.id
		LEFT JOIN
			statuses stat ON stat.id = p.statusId
		LEFT JOIN
			segments seg ON seg.id = p.segmentId
		WHERE 
			p.id = $1 AND p.deletedAt IS NULL
	`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	project := entities.Project{}
	userMap := make(map[int]bool)

	for rows.Next() {
		var dateStarted, dateDeadline *time.Time
		user := entities.User{}
		var userID, userAge, userDeletedBy, projectSegmentId, projectStatusId, userRoleId *int
		var userFirstName, userLastName, userEmail *string
		var userLastActiveAt, userCreatedAt, userUpdatedAt, userDeletedAt *time.Time

		var statusID, segmentId *int
		var statusName, statusDescription *string
		var segmentName, segmentDescription *string

		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Url, &project.Progress, &projectSegmentId, &projectStatusId, &dateStarted, &dateDeadline, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt, &project.DeletedBy,
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

		if projectSegmentId != nil {
			project.SegmentID = *segmentId
			project.Segment = entities.Segment{
				ID:          *projectSegmentId,
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

	return &project, nil
}

// CREATE PROJECT WITH USERS into users_projects
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
	err = tx.QueryRow("INSERT INTO projects (name, description, progress, url, statusId, segmentId, dateStarted, dateDeadline, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, name, description, progress, url, statusId, segmentId, dateStarted, dateDeadline, createdAt, updatedAt",
		payload.Name,
		payload.Description,
		progress,
		payload.Url,
		payload.StatusID,
		payload.SegmentID,
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
		&proj.SegmentID,
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

// UPDATE PROJECT WITH USERS into users_projects
func (s *Store) ProjectUpdate(payload entities.ProjectUpdatePayload) (map[string]interface{}, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	// Fetch the existing project
	query := `
		SELECT id, name, description, progress, url, dateStarted, dateDeadline, statusId, segmentId, createdAt, updatedAt
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
			&existProj.ID, &existProj.Name, &existProj.Description, &existProj.Progress, &existProj.Url,
			&existProj.DateStarted, &existProj.DateDeadline, &existProj.StatusID, &existProj.SegmentID,
			&existProj.CreatedAt, &existProj.UpdatedAt,
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

	if payload.Url != nil {
		existProj.Url = payload.Url
	}

	// Ensure non-nullable fields are not nullified
	if payload.StatusID != 0 {
		existProj.StatusID = payload.StatusID
	}
	if payload.SegmentID != 0 {
		existProj.SegmentID = payload.SegmentID
	}

	updateQuery := `
		UPDATE projects
		SET name = $1, description = $2, progress = $3, url = $4, dateStarted = $5, dateDeadline = $6, statusId = $7, segmentId = $8, updatedAt = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING id, name, description, progress, url, dateStarted, dateDeadline, statusId, segmentId, createdAt, updatedAt`
	proj := entities.Project{}
	err = tx.QueryRow(updateQuery,
		existProj.Name,
		existProj.Description,
		existProj.Progress,
		existProj.Url,
		existProj.DateStarted,
		existProj.DateDeadline,
		existProj.StatusID,
		existProj.SegmentID,
		existProj.ID,
	).Scan(
		&proj.ID, &proj.Name, &proj.Description, &proj.Progress, &proj.Url,
		&proj.DateStarted, &proj.DateDeadline, &proj.StatusID, &proj.SegmentID,
		&proj.CreatedAt, &proj.UpdatedAt,
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
		"segmentId":    proj.SegmentID,
		"url":          proj.Url,
		"dateStarted":  formatDate(proj.DateStarted),
		"dateDeadline": formatDate(proj.DateDeadline),
		"createdAt":    proj.CreatedAt,
		"updatedAt":    proj.UpdatedAt,
	}
}
