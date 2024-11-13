package role

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

func (s *Store) GetRoles() ([]*entities.Role, error) {
	rows, err := s.db.Query(`SELECT 
		name,
		description
	FROM roles`)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %v", err)
	}
	defer rows.Close()

	var roles []*entities.Role

	for rows.Next() {
		var role entities.Role

		err := rows.Scan(&role.Name, &role.Description)
		if err != nil {
			log.Printf("Failed to scan role: %v", err)
			continue
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over role rows: %v", err)
	}
	return roles, nil
}

func (s *Store) GetRole(id int) (*entities.Role, error) {
	rows, err := s.db.Query("Select id, name, description FROM roles WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	role := &entities.Role{}
	for rows.Next() {
		err := scanRowIntoRole(rows, role)
		if err != nil {
			return nil, err
		}
	}

	if role.ID == 0 {
		return nil, fmt.Errorf("role not found")
	}
	return role, nil
}

func (s *Store) CreateRole(role entities.Role) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO roles (name, description) VALUES (?, ?)", role.Name, role.Description)
	if err != nil {
		// If there's an error, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	// Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (s *Store) UpdateRole(role entities.Role) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	existingRole, err := s.GetRole(role.ID)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error fetching role: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	_, err = tx.Exec("UPDATE roles SET name = ?, description = ? WHERE id = ?", role.Name, role.Description, existingRole.ID)
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

	existingRole, err := s.GetRole(id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error fetching role: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	_, err = tx.Exec("DELETE FROM roles WHERE id = ?", existingRole.ID)
	if err != nil {
		// Rollback in case of any error during deletion
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
	)
}
