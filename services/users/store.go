package users

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
	"github.com/norrico31/it210-core-service-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Login(payload entities.UserLoginPayload) (string, entities.User, error) {
	user := entities.User{}
	err := s.db.QueryRow(`
		SELECT
			id, firstName, lastName, email, password, age, lastActiveAt, createdAt, updatedAt, deletedAt
		FROM users 
		WHERE email = $1 AND deletedAt IS NULL
	`, payload.Email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password,
		&user.Age, &user.LastActiveAt, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return "", user, fmt.Errorf("user not found")
	} else if err != nil {
		return "", user, fmt.Errorf("failed to query user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return "", user, fmt.Errorf("invalid password")
	}
	fmt.Println(user.ID)
	token, err := utils.GenerateJWT(user)
	if err != nil {
		return "", user, fmt.Errorf("failed to generate token: %v", err)
	}

	_, err = s.db.Exec(`UPDATE users SET lastActiveAt = NULL WHERE id = $1`, user.ID)
	if err != nil {
		log.Printf("Failed to update last active timestamp for user %d: %v", user.ID, err)
	}

	return token, user, nil
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
			p.statusId AS project_status_id,
			p.progress AS project_progress,
            p.description AS project_description,
            p.createdAt AS project_createdAt,
            p.updatedAt AS project_updatedAt,
			
			s.id AS status_id,
			s.name AS status_name,
			s.description AS status_description,
            s.createdAt AS status_createdAt,
            s.updatedAt AS status_updatedAt

        FROM users u
		LEFT JOIN roles r ON r.id = u.roleId
        LEFT JOIN users_projects up ON u.id = up.user_id
        LEFT JOIN projects p ON up.project_id = p.id
		LEFT JOIN statuses s ON s.id = p.statusId
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
			userID                                     int
			projectID                                  sql.NullInt32
			projectName                                sql.NullString
			projectDescription                         sql.NullString
			Progress                                   *float64
			projectCreatedAt                           sql.NullTime
			projectUpdatedAt                           sql.NullTime
			user                                       entities.User
			project                                    entities.Project
			userAge, roleId, statusID, projectStatusID *int
			roleName, roleDescription                  *string
			roleCreatedAt, roleUpdatedAt               *time.Time

			statusName, statusDescription    *string
			statusCreatedAt, statusUpdatedAt *time.Time
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
			&projectStatusID,
			&Progress,
			&projectDescription,
			&projectCreatedAt,
			&projectUpdatedAt,

			&statusID,
			&statusName,
			&statusDescription,
			&statusCreatedAt,
			&statusUpdatedAt,
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
			project.Progress = Progress
			project.StatusID = *projectStatusID
			project.CreatedAt = projectCreatedAt.Time
			project.UpdatedAt = projectUpdatedAt.Time

			if projectStatusID != nil {
				project.Status = entities.Status{
					ID:          *statusID,
					Name:        *statusName,
					Description: *statusDescription,
					CreatedAt:   *statusCreatedAt,
					UpdatedAt:   *statusUpdatedAt,
				}
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
			p.statusId AS project_status_id,
			p.progress AS project_progress,
			p.segmentId AS project_segment_id,
			p.description AS project_description,
			p.createdAt AS project_created_at,
			p.updatedAt AS project_updated_at,

			s.id AS status_id,
			s.name AS status_name,
			s.description AS status_description,
			s.createdAt AS status_created_at,
			s.updatedAt AS status_updated_at,
			
			seg.id AS segment_id,
			seg.name AS segment_name,
			seg.description AS segment_description,
			seg.createdAt AS segment_created_at,
			seg.updatedAt AS segment_updated_at

		FROM users u
		LEFT JOIN roles r ON r.id = u.roleId
		LEFT JOIN users_projects up ON u.id = up.user_id
		LEFT JOIN projects p ON up.project_id = p.id
		LEFT JOIN statuses s ON s.id = p.statusId
		LEFT JOIN segments seg ON seg.id = p.segmentId
		WHERE u.id = $1 AND u.deletedAt IS NULL
		ORDER BY p.id;
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by id: %v", err)
	}
	defer rows.Close()

	var (
		userID                                                                                                                                 int
		user                                                                                                                                   entities.User
		roleID, userAge, projectStatusID, statusID, projectSegmentID                                                                           *int
		roleName, roleDescription, statusName, statusDescription                                                                               *string
		roleCreatedAt, roleUpdatedAt, projectCreatedAt, projectUpdatedAt, statusCreatedAt, statusUpdatedAt, segmentCreatedAt, segmentUpdatedAt *time.Time
		project                                                                                                                                entities.Project
		segmentID                                                                                                                              *int
		segmentName, segmentDescription                                                                                                        *string
	)

	// Using a map to track the user and its associated projects
	userMap := make(map[int]*entities.User)

	for rows.Next() {
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

			&roleID,
			&roleName,
			&roleDescription,
			&roleCreatedAt,
			&roleUpdatedAt,

			&project.ID,
			&project.Name,
			&projectStatusID,
			&project.Progress,
			&projectSegmentID,
			&project.Description,
			&projectCreatedAt,
			&projectUpdatedAt,

			&statusID,
			&statusName,
			&statusDescription,
			&statusCreatedAt,
			&statusUpdatedAt,

			&segmentID,
			&segmentName,
			&segmentDescription,
			&segmentCreatedAt,
			&segmentUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user or project: %v", err)
		}

		// Set user attributes
		user.ID = userID
		if userAge != nil {
			user.Age = *userAge
		}
		if roleID != nil {
			user.Role = entities.Role{
				ID:          *roleID,
				Name:        *roleName,
				Description: *roleDescription,
				CreatedAt:   *roleCreatedAt,
				UpdatedAt:   *roleUpdatedAt,
			}
		}

		// Initialize user if not already done
		if _, exists := userMap[userID]; !exists {
			userMap[userID] = &user
		}

		// Add project if exists
		if project.ID != 0 {
			project.StatusID = *projectStatusID
			if statusID != nil {
				project.Status = entities.Status{
					ID:          *statusID,
					Name:        *statusName,
					Description: *statusDescription,
					CreatedAt:   *statusCreatedAt,
					UpdatedAt:   *statusUpdatedAt,
				}
			}
			project.CreatedAt = *projectCreatedAt
			project.UpdatedAt = *projectUpdatedAt

			// Add segment to project
			// if segmentID != nil {
			// 	project.SegmentID = *segmentID
			// 	project.Segment = entities.Segment{
			// 		ID:          *segmentID,
			// 		Name:        *segmentName,
			// 		Description: *segmentDescription,
			// 		CreatedAt:   *segmentCreatedAt,
			// 		UpdatedAt:   *segmentUpdatedAt,
			// 	}
			// }

			// Append project to user
			userMap[userID].Projects = append(userMap[userID].Projects, project)
		}
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", rows.Err())
	}

	// Return the first user found (there should only be one)
	if len(userMap) > 0 {
		return userMap[userID], nil
	}

	return nil, fmt.Errorf("user with id %d not found", id)
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

// TODO: ADD USERID IN SEGMENT and OTHER RELATIONS HERE
func (s *Store) UpdateUser(userId int, user entities.UserUpdatePayload, projectIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

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

	_, err = tx.Exec(`DELETE FROM users_projects WHERE user_id = $1`, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear project associations: %v", err)
	}

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

	//* OPTION A
	// CHECK IF INPUT PROJECT EXIST IN USERS AND GET IT TO SET THE deletedAt to null
	// IF PROJECT DOESNS'T EXIST INSERT THE NEW PROJECTS in users_projects
	//* OPTION B CREATE NEW TABLE FOR users_projects_deleted for softdelete

	// Delete existing in users_projects
	// _, err = tx.Exec(`UPDATE users_projects SET deletedAt = CURRENT_TIMESTAMP, deletedBy = $1 WHERE user_id = $1`, userId)
	// if err != nil {
	// 	tx.Rollback()
	// 	return fmt.Errorf("failed to clear project associations: %v", err)
	// }

	// // Create new user with projects in users_projects
	// for _, projID := range projectIDs {
	// 	// _, err = tx.Exec(`
	// 	// 	INSERT INTO users_projects (user_id, project_id)
	// 	// 	VALUES ($1, $2)
	// 	// `, userId, projID)
	// 	_, err = tx.Exec(`
	// 		UPDATE users_projects SET deletedAt = NULL, deletedBy = NULL WHERE user_id = $1
	// 	`, userId)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return fmt.Errorf("failed to associate user with project %d: %v", projID, err)
	// 	}
	// }

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %v", err)
	}

	return nil
}

func (s *Store) DeleteUser(userId int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// TODO: ADD THE USER WHO DELETED DATA
	_, err = tx.Exec("UPDATE users SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1", userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user with project %d: %v", userId, err)
	}
	// TODO: ADD THE USER WHO DELETED DATA

	_, err = tx.Exec("UPDATE users_projects SET deletedAt = CURRENT_TIMESTAMP WHERE user_id = $1", userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user with project %d: %v", userId, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %v", err)
	}
	return nil
}

func (s *Store) RestoreUser(userId int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// TODO: ADD THE USER WHO DELETED DATA
	_, err = tx.Exec("UPDATE users SET deletedAt = NULL WHERE id = $1", userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to restore user with project %d: %v", userId, err)
	}

	// TODO: ADD THE USER WHO DELETED DATA
	_, err = tx.Exec("UPDATE users_projects SET deletedAt = NULL WHERE user_id = $1", userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to restore user with project %d: %v", userId, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %v", err)
	}
	return nil
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
