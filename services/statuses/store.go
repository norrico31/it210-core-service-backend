package statuses

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

func (s *Store) GetStatuses() ([]entities.Status, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, createdAt, updatedAt
			FROM statuses
		WHERE deletedAt IS NULL
		ORDER BY createdAt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query statuses: %v", err)
	}
	defer rows.Close()

	statuses := []entities.Status{}

	for rows.Next() {
		statuse := entities.Status{}

		err := scanRowIntoStatus(rows, &statuse)
		if err != nil {
			log.Printf("Failed to scan statuse: %v", err)
			continue
		}
		statuses = append(statuses, statuse)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over statuse rows: %v", err)
	}
	return statuses, nil
}

func (s *Store) GetStatus(id int) (*entities.Status, error) {
	statuse := entities.Status{}
	err := s.db.QueryRow("SELECT id, name, description, createdAt, updatedAt FROM statuses WHERE deletedAt IS NULL AND id = $1", id).Scan(
		&statuse.ID,
		&statuse.Name,
		&statuse.Description,
		&statuse.CreatedAt,
		&statuse.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("statuse not found")
	}

	if statuse.ID == 0 {
		return nil, fmt.Errorf("statuse not found")
	}
	return &statuse, nil
}

func (s *Store) CreateStatus(payload entities.StatusPayload) (*entities.Status, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO statuses (name, description) VALUES ($1, $2)", payload.Name, payload.Description)
	if err != nil {
		// If there's an error, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return nil, err
	}

	// Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return &entities.Status{Name: payload.Name, Description: payload.Description}, err
	}

	return &entities.Status{Name: payload.Name, Description: payload.Description}, err
}

func (s *Store) UpdateStatus(payload entities.StatusPayload) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE statuses SET name = $1, description = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3", payload.Name, payload.Description, payload.ID)
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

func (s *Store) DeleteStatus(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE statuses SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1", id)

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

func (s *Store) RestoreStatus(id int) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE statuses SET deletedAt = NULL WHERE id = $1", id)
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

func scanRowIntoStatus(rows *sql.Rows, statuse *entities.Status) error {
	return rows.Scan(
		&statuse.ID,
		&statuse.Name,
		&statuse.Description,
		&statuse.CreatedAt,
		&statuse.UpdatedAt,
	)
}
