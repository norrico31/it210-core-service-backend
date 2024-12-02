package priorities

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

func (s *Store) GetPriorities() ([]entities.Priority, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, createdAt, updatedAt
			FROM priorities
		WHERE deletedAt IS NULL
		ORDER BY createdAt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query priorities: %v", err)
	}
	defer rows.Close()

	priorities := []entities.Priority{}

	for rows.Next() {
		priority := entities.Priority{}

		err := scanRowIntoPriority(rows, &priority)
		if err != nil {
			log.Printf("Failed to scan priority: %v", err)
			continue
		}
		priorities = append(priorities, priority)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over priority rows: %v", err)
	}
	return priorities, nil
}

func (s *Store) GetPriority(id int) (*entities.Priority, error) {
	priority := entities.Priority{}
	err := s.db.QueryRow("SELECT id, name, description, createdAt, updatedAt FROM priorities WHERE deletedAt IS NULL AND id = $1", id).Scan(
		&priority.ID,
		&priority.Name,
		&priority.Description,
		&priority.CreatedAt,
		&priority.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("priority not found")
	}

	if priority.ID == 0 {
		return nil, fmt.Errorf("priority not found")
	}
	return &priority, nil
}

func (s *Store) CreatePriority(payload entities.PriorityPayload) (*entities.Priority, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO priorities (name, description) VALUES ($1, $2)", payload.Name, payload.Description)
	if err != nil {
		// If there's an error, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return nil, err
	}

	// Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return &entities.Priority{Name: payload.Name, Description: payload.Description}, err
	}

	return &entities.Priority{Name: payload.Name, Description: payload.Description}, err
}

func (s *Store) UpdatePriority(payload entities.PriorityPayload) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE priorities SET name = $1, description = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3", payload.Name, payload.Description, payload.ID)
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

func (s *Store) DeletePriority(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE priorities SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1", id)

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

func (s *Store) RestorePriority(id int) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE priorities SET deletedAt = NULL WHERE id = $1", id)
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

func scanRowIntoPriority(rows *sql.Rows, priority *entities.Priority) error {
	return rows.Scan(
		&priority.ID,
		&priority.Name,
		&priority.Description,
		&priority.CreatedAt,
		&priority.UpdatedAt,
	)
}
