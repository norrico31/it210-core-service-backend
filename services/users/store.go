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
	rows, err := s.db.Query(`SELECT 
		id,
		firstName,
		lastName,
		email,
		age,
		lastActiveAt,
		createdAt,
		updatedAt,
		deletedAt
	FROM users WHERE deletedAt IS NULL`)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	var users []*entities.User

	fmt.Println("MAY SOMETHING DITO GET USERS")
	for rows.Next() {
		var user entities.User
		err := scanRowIntoUser(rows, &user)
		if err != nil {
			log.Printf("Failed to scan user: %v", err)
			continue
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over user rows: %v", err)
	}
	return users, nil
}

func (s *Store) GetUserById(id int) (*entities.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(entities.User)
	for rows.Next() {
		err := scanRowIntoUser(rows, user)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
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
