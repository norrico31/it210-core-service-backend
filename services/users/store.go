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
            u.lastActiveAt,
            u.createdAt,
            u.updatedAt,
            u.deletedAt,
            p.id AS project_id,
            p.name AS project_name,
            p.description AS project_description,
            p.createdAt AS project_createdAt,
            p.updatedAt AS project_updatedAt,
            p.deletedAt AS project_deletedAt
        FROM users u
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
			userID             int
			projectID          sql.NullInt32
			projectName        sql.NullString
			projectDescription sql.NullString
			projectCreatedAt   sql.NullTime
			projectUpdatedAt   sql.NullTime
			projectDeletedAt   sql.NullTime
			user               entities.User
			project            entities.Project
		)

		// Scan user and project data from the row
		err := rows.Scan(
			&userID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Age,
			&user.LastActiveAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&projectID,
			&projectName,
			&projectDescription,
			&projectCreatedAt,
			&projectUpdatedAt,
			&projectDeletedAt,
		)
		if err != nil {
			log.Printf("Failed to scan user or project: %v", err)
			continue
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

			if projectDeletedAt.Valid {
				project.DeletedAt = &projectDeletedAt.Time
			} else {
				project.DeletedAt = nil
			}

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
	err := s.db.QueryRow(`
		SELECT  
			id, firstName,
			age, lastName, 
			e-mail, lastActiveAt,
			roleId,
			createdAt, updatedAt, 
		FROM users 
		WHERE deletedAt IS NULL id = $1
	`, id).Scan()
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

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
		fmt.Println("DITO NGA YON", err)
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

func (s *Store) CreateUser(user entities.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password, lastActiveAt) VALUES (?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password, nil)
	return err
}

func (s *Store) UpdateUser(user entities.User) error {
	_, err := s.db.Exec("UPDATE users SET firstName = ?, lastName = ?, email = ?, password = ? WHERE id = ?", user.FirstName, user.LastName, user.Email, user.Password, user.ID)
	return err
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
