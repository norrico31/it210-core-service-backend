package users

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUsers() ([]*entities.User, error) {
	rows, err := s.db.Query(`
        SELECT 
            u.id AS user_id,
            u.firstName,
            u.lastName,
            u.email,
            u.age,
			u.roleId,
            u.lastActiveAt,
            u.createdAt,
            u.updatedAt,
			u.deletedAt,

			r.id AS role_id,
			r.name AS role_name,
			r.description AS role_description,
			r.createdAt AS role_created_at,
			r.updatedAt AS role_updated_at,

            p.id AS project_id,
            p.name AS project_name,
            p.description AS project_description,
            p.createdAt AS project_createdAt,
            p.updatedAt AS project_updatedAt

        FROM users u
		LEFT JOIN roles r ON r.id = u.roleId
        LEFT JOIN users_projects up ON u.id = up.user_id
        LEFT JOIN projects p ON up.project_id = p.id
        WHERE u.deletedAt IS NULL
        ORDER BY u.id, p.id;
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to query users and projects: %v", err)
	}
	defer rows.Close()

	userMap := make(map[int]*entities.User)

	for rows.Next() {
		var (
			userID                       int
			projectID                    sql.NullInt32
			projectName                  sql.NullString
			projectDescription           sql.NullString
			projectCreatedAt             sql.NullTime
			projectUpdatedAt             sql.NullTime
			user                         entities.User
			project                      entities.Project
			userAge, roleId              *int
			roleName, roleDescription    *string
			roleCreatedAt, roleUpdatedAt *time.Time
		)

		err := rows.Scan(
			&userID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&userAge,
			&user.RoleId,
			&user.LastActiveAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,

			&roleId,
			&roleName,
			&roleDescription,
			&roleCreatedAt,
			&roleUpdatedAt,

			&projectID,
			&projectName,
			&projectDescription,
			&projectCreatedAt,
			&projectUpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan user or project: %v", err)
			continue
		}
		if userAge != nil {
			user.Age = *userAge
		}
		if roleId != nil {
			user.Role = entities.Role{
				ID:          *roleId,
				Name:        *roleName,
				Description: *roleDescription,
				CreatedAt:   *roleCreatedAt,
				UpdatedAt:   *roleUpdatedAt,
			}
		}

		// Check if the user already exists in the map
		if _, exists := userMap[userID]; !exists {
			user.ID = userID
			user.Projects = []entities.Project{}
			userMap[userID] = &user
		}

		// If project information is valid, add it to the user's projects
		if projectID.Valid {
			project.ID = int(projectID.Int32)
			project.Name = projectName.String
			project.Description = projectDescription.String
			project.CreatedAt = projectCreatedAt.Time
			project.UpdatedAt = projectUpdatedAt.Time

			// Append the project to the user's project list
			userMap[userID].Projects = append(userMap[userID].Projects, project)
		}
	}

	// Convert userMap to a slice
	userList := make([]*entities.User, 0, len(userMap))
	for _, user := range userMap {
		userList = append(userList, user)
	}

	return userList, nil
}

func (s *Store) GetUserById(id int) (*entities.User, error) {
	user := entities.User{}
	var (
		roleID, userAge              *int
		roleName, roleDescription    *string
		roleCreatedAt, roleUpdatedAt *time.Time
	)

	err := s.db.QueryRow(`
		SELECT 
            u.id AS user_id,
            u.firstName,
            u.lastName,
            u.email,
            u.age,
			u.roleId,
            u.lastActiveAt,
            u.createdAt,
            u.updatedAt,
            u.deletedAt,

			r.id AS role_id,
			r.name AS role_name,
			r.description AS role_description,
			r.createdAt AS role_created_at,
			r.updatedAt AS role_updated_at

        FROM users u
		LEFT JOIN roles r ON r.deletedAt IS NULL AND r.id = u.roleId
		WHERE u.deletedAt IS NULL AND u.id = $1
	`, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&userAge,
		&user.RoleId,
		&user.LastActiveAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&roleID,
		&roleName,
		&roleDescription,
		&roleCreatedAt,
		&roleUpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user id %v not found!", id)
	}

	if userAge != nil {
		user.Age = *userAge
	}

	if user.RoleId != nil {
		user.Role = entities.Role{
			ID:          *roleID,
			Name:        *roleName,
			Description: *roleDescription,
			CreatedAt:   *roleCreatedAt,
			UpdatedAt:   *roleUpdatedAt,
		}
	}

	rows, err := s.db.Query(`
		SELECT 
			p.id AS project_id,
			p.name AS project_name,
			p.description AS project_description,
			p.progress,
			p.dateStarted,
			p.dateDeadline,
			p.createdAt AS project_created_at,
			p.updatedAt AS project_updated_at

		FROM users_projects up
		LEFT JOIN projects p ON up.project_id = p.id
		WHERE up.user_id = $1 AND p.deletedAt IS NULL
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query users_projects: %v", err)
	}
	defer rows.Close()

	projects := []entities.Project{}
	for rows.Next() {
		proj := entities.Project{}
		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan project: %v", err)
		}
		projects = append(projects, proj)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating over projects rows: %v", rows.Err())
	}

	user.Projects = projects

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func (s *Store) GetUserByEmail(email string) (*entities.User, error) {
	// SQL query to fetch all relevant fields
	query := `
        SELECT id, firstName, age, lastName, email, password, lastActiveAt, createdAt, updatedAt, deletedAt
        FROM users WHERE email = $1`
	row := s.db.QueryRow(query, email)

	user := &entities.User{}
	var lastActiveAt sql.NullTime
	var deletedAt sql.NullTime

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.Age,
		&user.LastName,
		&user.Email,
		&user.Password,
		&lastActiveAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		// Check if the error is because no rows were found
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}

	// Handle nullable fields
	if lastActiveAt.Valid {
		user.LastActiveAt = &lastActiveAt.Time
	} else {
		user.LastActiveAt = nil
	}

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	} else {
		user.DeletedAt = nil
	}
	// Return the populated user struct
	return user, nil
}

type ProjectValidate struct {
	ID   string `json:"name"`
	Name string `json:"description"`
}

func (s *Store) CreateUser(payload entities.UserCreatePayload) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	var userID int
	err = tx.QueryRow(`
		INSERT INTO users (firstName, lastName, email, roleId, age, password, createdAt, updatedAt) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		payload.FirstName,
		payload.LastName,
		payload.Email,
		payload.RoleId,
		payload.Age,
		payload.Password,
		time.Now(),
		time.Now(),
	).Scan(&userID)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	// Insert user-project associations if any
	if payload.ProjectIDS != nil {
		for _, projID := range *payload.ProjectIDS {
			_, err := tx.Exec(`
				INSERT INTO users_projects (user_id, project_id)
				VALUES ($1, $2)
			`,
				userID,
				projID)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					return fmt.Errorf("association insert error: %v, rollback error: %v", err, rbErr)
				}
				return err
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %v", err)
	}

	return nil
}

func (s *Store) UpdateUser(userId int, user entities.UserUpdatePayload, projectIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Update user
	_, err = tx.Exec(`
		UPDATE users
		SET firstName = $1, lastName = $2, email = $3, roleId = $4, age = $5, password = $6, updatedAt = $7
		WHERE id = $8
	`,
		user.FirstName,
		user.LastName,
		user.Email,
		user.RoleId,
		user.Age,
		user.Password,
		time.Now(),
		userId,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user: %v", err)
	}

	// Delete existing in users_projects
	_, err = tx.Exec(`DELETE FROM users_projects WHERE user_id = $1`, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear project associations: %v", err)
	}

	// Create new user with projects in users_projects
	for _, projID := range projectIDs {
		_, err = tx.Exec(`
			INSERT INTO users_projects (user_id, project_id)
			VALUES ($1, $2)
		`, userId, projID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to associate user with project %d: %v", projID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %v", err)
	}

	return nil
}

func (s *Store) DeleteUser(id int) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (s *Store) SetUserActive(userId int) error {
	result, err := s.db.Exec("UPDATE users SET lastActiveAt = NOW() WHERE id = $1", userId)
	if err != nil {
		fmt.Println("Error updating user active status:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Rows affected: %d\n", rowsAffected)

	return nil
}
func (s *Store) UpdateLastActiveTime(userId int, time time.Time) error {
	_, err := s.db.Exec("UPDATE users SET lastActiveAt = ? WHERE id = ?", time, userId)
	return err
}

func scanRowIntoUser(rows *sql.Rows, user *entities.User) error {
	return rows.Scan(
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
}
