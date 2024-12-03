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
	// Query segments and associated projects, including segments without projects
	rows, err := s.db.Query(`
		SELECT 
			seg.id AS segment_id, 
			seg.name AS segment_name, 
			seg.description AS segment_description,
			seg.createdAt AS segment_createdAt,
			seg.updatedAt AS segment_updatedAt,
			p.id AS project_id, 
			p.name AS project_name, 
			p.description AS project_description,
			p.progress AS project_progress, 
			p.url AS project_url, 
			p.dateStarted AS project_dateStarted, 
			p.dateDeadline AS project_dateDeadline, 
			p.createdAt AS project_createdAt,
			p.updatedAt AS project_updatedAt
		FROM segments seg
		LEFT JOIN segments_projects sp ON seg.id = sp.segmentId
		LEFT JOIN projects p ON sp.projectId = p.id
		ORDER BY seg.id, p.createdAt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query segments and projects: %v", err)
	}
	defer rows.Close()

	// Create a map to associate segment ID with the segment itself
	segments := make(map[int]*entities.Segment)

	// Process each row
	for rows.Next() {
		var segment entities.Segment
		var project entities.Project

		// Use sql.NullInt64 for project.ID to handle NULL values
		var projectID sql.NullInt64

		err := rows.Scan(
			&segment.ID, &segment.Name, &segment.Description, &segment.CreatedAt, &segment.UpdatedAt,
			&projectID, &project.Name, &project.Description, &project.Progress, &project.Url,
			&project.DateStarted, &project.DateDeadline, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		// Check if the segment is already in the map
		if _, exists := segments[segment.ID]; !exists {
			segments[segment.ID] = &segment
		}

		// Only add the project if the project ID is not NULL
		if projectID.Valid {
			project.ID = int(projectID.Int64) // Convert int64 to int
			segments[segment.ID].Projects = append(segments[segment.ID].Projects, project)
		}
	}

	// Convert the map back to a slice
	var result []entities.Segment
	for _, segment := range segments {
		result = append(result, *segment)
	}

	// Check for any row iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return result, nil
}

func (s *Store) GetSegment(id int) (*entities.Segment, error) {
	// Query segment and its associated projects
	rows, err := s.db.Query(`
		SELECT 
			seg.id AS segment_id, 
			seg.name AS segment_name, 
			seg.description AS segment_description,
			seg.createdAt AS segment_createdAt,
			seg.updatedAt AS segment_updatedAt,
			p.id AS project_id, 
			p.name AS project_name, 
			p.description AS project_description,
			p.progress AS project_progress, 
			p.url AS project_url, 
			p.dateStarted AS project_dateStarted, 
			p.dateDeadline AS project_dateDeadline, 
			p.createdAt AS project_createdAt,
			p.updatedAt AS project_updatedAt
		FROM segments seg
		LEFT JOIN segments_projects sp ON seg.id = sp.segmentId
		LEFT JOIN projects p ON sp.projectId = p.id
		WHERE seg.id = $1 AND seg.deletedAt IS NULL
		ORDER BY p.createdAt DESC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query segment and projects: %v", err)
	}
	defer rows.Close()

	// Create a map to associate segment ID with the segment itself
	segment := &entities.Segment{}
	// Process each row
	for rows.Next() {
		var project entities.Project
		var projectID sql.NullInt64

		err := rows.Scan(
			&segment.ID, &segment.Name, &segment.Description, &segment.CreatedAt, &segment.UpdatedAt,
			&projectID, &project.Name, &project.Description, &project.Progress, &project.Url,
			&project.DateStarted, &project.DateDeadline, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		// Only add the project if the project ID is not NULL
		if projectID.Valid {
			project.ID = int(projectID.Int64) // Convert int64 to int
			segment.Projects = append(segment.Projects, project)
		}
	}

	// Check for any row iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	if segment.ID == 0 {
		return nil, fmt.Errorf("segment not found")
	}

	return segment, nil
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
