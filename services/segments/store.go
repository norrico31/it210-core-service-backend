package segments

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

func (s *Store) GetSegments() ([]entities.Segment, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, createdAt, updatedAt
			FROM segments
		WHERE deletedAt IS NULL
		ORDER BY createdAt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query segments: %v", err)
	}
	defer rows.Close()

	segments := []entities.Segment{}

	for rows.Next() {
		segment := entities.Segment{}

		err := scanRowIntoSegment(rows, &segment)
		if err != nil {
			log.Printf("Failed to scan segment: %v", err)
			continue
		}
		segments = append(segments, segment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over segment rows: %v", err)
	}
	return segments, nil
}

func (s *Store) GetSegment(id int) (*entities.Segment, error) {
	segment := entities.Segment{}
	err := s.db.QueryRow("SELECT id, name, description, createdAt, updatedAt FROM segments WHERE deletedAt IS NULL AND id = $1", id).Scan(
		&segment.ID,
		&segment.Name,
		&segment.Description,
		&segment.CreatedAt,
		&segment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("segment not found")
	}

	if segment.ID == 0 {
		return nil, fmt.Errorf("segment not found")
	}
	return &segment, nil
}

func (s *Store) CreateSegment(payload entities.SegmentPayload) (*entities.Segment, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO segments (name, description) VALUES ($1, $2)", payload.Name, payload.Description)
	if err != nil {
		// If there's an error, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return nil, err
	}

	// Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return &entities.Segment{Name: payload.Name, Description: payload.Description}, err
	}

	return &entities.Segment{Name: payload.Name, Description: payload.Description}, err
}

func (s *Store) UpdateSegment(payload entities.SegmentPayload) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE segments SET name = $1, description = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3", payload.Name, payload.Description, payload.ID)
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

func (s *Store) DeleteSegment(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE segments SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1", id)

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

func (s *Store) RestoreSegment(id int) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE segments SET deletedAt = NULL WHERE id = $1", id)
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

func scanRowIntoSegment(rows *sql.Rows, segment *entities.Segment) error {
	return rows.Scan(
		&segment.ID,
		&segment.Name,
		&segment.Description,
		&segment.CreatedAt,
		&segment.UpdatedAt,
	)
}
