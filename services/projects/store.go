package projects

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

func (s *Store) GetProjects() ([]*entities.Project, error) {
	rows, err := s.db.Query(`
		SELECT
			p.id AS project_id,
			p.name AS project_name,
			p.description AS project_description,
			p.createdAt AS project_created_at,
			p.updatedAt AS project_updated_at,
			p.deletedAt AS project_deleted_at,
			u.id AS user_id,
			u.firstName AS user_first_name,
			u.lastName AS user_last_name,
			u.email AS user_email,
			u.age AS user_age,
			u.lastActiveAt AS user_last_active_at,
			u.createdAt AS user_created_at,
			u.updatedAt AS user_updated_at,
			u.deletedAt AS user_deleted_at
		FROM
			projects p
		LEFT JOIN
			users_projects up ON p.id = up.project_id
		LEFT JOIN
			users u ON up.user_id = u.id;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projectsMap := make(map[int]*entities.Project)
	for rows.Next() {
		var projectID int
		var project entities.Project
		var user entities.User

		// Scan project data
		err = rows.Scan(
			&projectID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
			&project.UpdatedAt,
			&project.DeletedAt,
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Age,
			&user.LastActiveAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Add project to map if not already present
		if _, exists := projectsMap[projectID]; !exists {
			project.ID = projectID
			project.Users = []entities.User{}
			projectsMap[projectID] = &project
		}

		// Add user to the project's users list if user data exists
		if user.ID != 0 {
			projectsMap[projectID].Users = append(projectsMap[projectID].Users, user)
		}
	}

	// Convert map to slice
	projects := []*entities.Project{}
	for _, project := range projectsMap {
		projects = append(projects, project)
	}

	return projects, nil
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
		fmt.Println("aha?")
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
