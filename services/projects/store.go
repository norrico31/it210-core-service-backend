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

// TODO: WHEN DELETE PLEASE INCLUDE THE USER ID TO TRACK WHO'S USER DELETE THE OPERATIONS (deletedBy in tables)
func (s *Store) GetProjects(str string) ([]*entities.Project, error) {
	query := fmt.Sprintf(`SELECT
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
		WHERE p.deletedAt %s
		`, str)
	rows, err := s.db.Query(query)

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
	query := `
		SELECT 
			p.id, p.name, p.description, p.createdAt, p.updatedAt, p.deletedAt,

			t.id AS task_id, t.title AS task_title, t.subTask AS task_subTask, 
			t.description AS task_description, t.statusId AS task_statusId, 
			t.userId AS task_userId, t.createdAt AS task_createdAt, 
			t.updatedAt AS task_updatedAt, t.deletedAt AS task_deletedAt,

			s.id AS status_id, s.name AS status_name, s.description AS status_description, 
			s.createdAt AS status_createdAt, s.updatedAt AS status_updatedAt, s.deletedAt AS status_deletedAt,

			u.id AS user_id, u.firstName AS user_firstName, u.lastName AS user_lastName, 
			u.email AS user_email, u.age AS user_age, u.lastActiveAt, 
			u.createdAt AS user_createdAt, u.updatedAt AS user_updatedAt, u.deletedAt AS user_deletedAt,

			up_user.id AS project_user_id, up_user.firstName AS project_user_firstName, 
			up_user.lastName AS project_user_lastName, up_user.email AS project_user_email, 
			up_user.age AS project_user_age, up_user.lastActiveAt AS project_user_lastActiveAt,
			up_user.createdAt AS project_user_createdAt, up_user.updatedAt AS project_user_updatedAt, 
			up_user.deletedAt AS project_user_deletedAt

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
			p.id = $1
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
			taskID, statusID, userID, projectUserID int
			task                                    entities.Task
			status                                  entities.Status
			user                                    entities.User
			projectUser                             entities.User
		)

		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt, &project.DeletedAt,

			&taskID, &task.Title, &task.SubTask, &task.Description, &task.StatusID,
			&task.UserID, &task.CreatedAt, &task.UpdatedAt, &task.DeletedAt,

			&statusID, &status.Name, &status.Description, &status.CreatedAt,
			&status.UpdatedAt, &status.DeletedAt,

			&userID, &user.FirstName, &user.LastName, &user.Email, &user.Age,
			&user.LastActiveAt, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,

			&projectUserID, &projectUser.FirstName, &projectUser.LastName,
			&projectUser.Email, &projectUser.Age, &projectUser.LastActiveAt,
			&projectUser.CreatedAt, &projectUser.UpdatedAt, &projectUser.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Add tasks
		if taskID != 0 {
			if _, exists := taskMap[taskID]; !exists {
				task.ID = taskID
				if statusID != 0 {
					status.ID = statusID
					task.Status = status
				}
				if userID != 0 {
					user.ID = userID
					task.User = user
				}
				taskMap[taskID] = task
			}
		}

		// Add users
		if projectUserID != 0 {
			if _, exists := userMap[projectUserID]; !exists {
				projectUser.ID = projectUserID
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
	fmt.Println("hala")
	proj := entities.Project{}
	err = tx.QueryRow("UPDATE projects SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id, name, description, createdAt, updatedAt, deletedAt", id).Scan(
		&proj.ID,
		&proj.Name,
		&proj.Description,
		&proj.CreatedAt,
		&proj.UpdatedAt,
		&proj.DeletedAt,
	)
	fmt.Println("hala?")

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error deleting: %v rollback error: %v", err, rollbackErr)
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return &proj, nil
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
	fmt.Println("hala?")
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
