package roles

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/norrico31/it210-core-service-backend/entities"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetRoles() ([]entities.Role, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, createdAt, updatedAt
			FROM roles
		WHERE deletedAt IS NULL
		ORDER BY createdAt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %v", err)
	}
	defer rows.Close()

	roles := []entities.Role{}

	for rows.Next() {
		role := entities.Role{}

		err := scanRowIntoRole(rows, &role)
		if err != nil {
			log.Printf("Failed to scan role: %v", err)
			continue
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over role rows: %v", err)
	}
	return roles, nil
}

func (s *Store) GetRole(id int) (*entities.Role, error) {
	role := entities.Role{}
	err := s.db.QueryRow("SELECT id, name, description, createdAt, updatedAt FROM roles WHERE deletedAt IS NULL AND id = $1", id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("role not found")
	}

	if role.ID == 0 {
		return nil, fmt.Errorf("role not found")
	}
	return &role, nil
}

func (s *Store) CreateRole(payload entities.RolePayload) (*entities.Role, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO roles (name, description) VALUES ($1, $2)", payload.Name, payload.Description)
	if err != nil {
		// If there's an error, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return nil, err
	}

	// Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return &entities.Role{Name: payload.Name, Description: payload.Description}, err
	}

	return &entities.Role{Name: payload.Name, Description: payload.Description}, err
}

func (s *Store) UpdateRole(payload entities.RolePayload) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE roles SET name = $1, description = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3", payload.Name, payload.Description, payload.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (s *Store) DeleteRole(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE roles SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1", id)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("delete error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (s *Store) RestoreRole(id int) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE roles SET deletedAt = NULL WHERE id = $1", id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("delete error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	// Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func scanRowIntoRole(rows *sql.Rows, role *entities.Role) error {
	return rows.Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
}
