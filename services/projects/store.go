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

func (s *Store) GetProjects() ([]*entities.Project, error) {
	rows, err := s.db.Query(
		`SELECT
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
			users u ON up.user_id = u.id
		WHERE p.deletedAt IS NULL
		`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to hold projects indexed by projectID
	projectsMap := make(map[int]*entities.Project)

	// Loop through query results
	for rows.Next() {
		var projectID int
		var project entities.Project
		var user entities.User

		// Use pointer to handle NULL values
		var userID *int
		var userFirstName, userLastName, userEmail *string
		var userAge *int
		var userLastActiveAt, userCreatedAt, userUpdatedAt, userDeletedAt *time.Time

		// Scan project and user data
		err = rows.Scan(
			&projectID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
			&project.UpdatedAt,
			&project.DeletedAt,
			&userID,
			&userFirstName,
			&userLastName,
			&userEmail,
			&userAge,
			&userLastActiveAt,
			&userCreatedAt,
			&userUpdatedAt,
			&userDeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// If project is not yet in map, add it
		if _, exists := projectsMap[projectID]; !exists {
			project.ID = projectID
			project.Users = []entities.User{} // Initialize empty users slice
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
			user.CreatedAt = *userCreatedAt
			user.UpdatedAt = *userUpdatedAt
			user.DeletedAt = userDeletedAt

			// Add user to the project's user list
			projectsMap[projectID].Users = append(projectsMap[projectID].Users, user)
		}
	}

	// Convert map to slice
	var projects []*entities.Project
	for _, project := range projectsMap {
		projects = append(projects, project)
	}

	//
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].CreatedAt.After(projects[j].CreatedAt)
	})
	return projects, nil
}

func (s *Store) GetProject(id int) (*entities.Project, error) {
	rows, err := s.db.Query("SELECT * FROM projects WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	proj := entities.Project{}
	for rows.Next() {
		err := scanRowIntoProject(rows, &proj)
		if err != nil {
			return nil, err
		}

	}

	if proj.ID == 0 {
		return nil, fmt.Errorf("project not found")
	}

	return &proj, nil
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
	fmt.Println("hala")
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
