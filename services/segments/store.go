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

// TODO GET WITH ASSOCIATE
func (s *Store) GetSegments() ([]entities.Segment, error) {
	// First, query to get all segments
	segmentRows, err := s.db.Query(`
		SELECT 
			seg.id, seg.name, seg.description, seg.createdAt, seg.updatedAt
		FROM segments seg
		ORDER BY seg.createdAt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query segments: %v", err)
	}
	defer segmentRows.Close()

	segments := []entities.Segment{}
	segmentMap := make(map[int]*entities.Segment) // Map to associate segments with their projects

	for segmentRows.Next() {
		var segment entities.Segment
		err := segmentRows.Scan(
			&segment.ID,
			&segment.Name,
			&segment.Description,
			&segment.CreatedAt,
			&segment.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan segment row: %v", err)
			continue
		}

		// Initialize the segment in the map
		segmentMap[segment.ID] = &segment
		segments = append(segments, segment)
	}

	// Check for any row iteration errors for segments
	if err := segmentRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over segment rows: %v", err)
	}

	// Next, query to get projects associated with the segments
	projectRows, err := s.db.Query(`
		SELECT 
			p.id as project_id, p.name as project_name, p.description as project_description,
			p.progress as project_progress, p.url as project_url, 
			p.dateStarted as project_dateStarted, p.dateDeadline as project_dateDeadline,
			p.createdAt as project_createdAt, p.updatedAt as project_updatedAt, 
			p.segmentId as project_segmentId
		FROM projects p
		ORDER BY p.createdAt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %v", err)
	}
	defer projectRows.Close()

	// Associate projects with their respective segments
	for projectRows.Next() {
		var project entities.Project
		var segmentId int

		err := projectRows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.Progress,
			&project.Url,
			&project.DateStarted,
			&project.DateDeadline,
			&project.CreatedAt,
			&project.UpdatedAt,
			&segmentId,
		)
		if err != nil {
			log.Printf("Failed to scan project row: %v", err)
			continue
		}

		// Check if the project is associated with the current segment (only append if IDs match)
		if segment, exists := segmentMap[segmentId]; exists && segment.ID == segmentId {
			// If the segment exists and IDs match, append the project to the segment's Projects slice
			segment.Projects = append(segment.Projects, project)
		}
	}

	// Check for any row iteration errors for projects
	if err := projectRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over project rows: %v", err)
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
