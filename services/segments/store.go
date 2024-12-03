package segments

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

func (s *Store) GetSegments() ([]entities.Segment, error) {
	rows, err := s.db.Query(`
		SELECT 
			seg.id AS segment_id, 
			seg.name AS segment_name, 
			seg.description AS segment_description,
			seg.createdAt AS segment_createdAt,
			seg.updatedAt AS segment_updatedAt,
			seg.deletedAt AS segment_deletedAt,
			p.id AS project_id, 
			p.name AS project_name, 
			p.description AS project_description,
			p.progress AS project_progress, 
			p.url AS project_url, 
			p.dateStarted AS project_dateStarted, 
			p.dateDeadline AS project_dateDeadline, 
			p.createdAt AS project_createdAt,
			p.updatedAt AS project_updatedAt,
			p.deletedAt AS project_deletedAt
		FROM segments seg
		LEFT JOIN segments_projects sp ON seg.id = sp.segmentId
		LEFT JOIN projects p ON sp.projectId = p.id
		WHERE seg.deletedAt IS NULL
		ORDER BY seg.id, p.createdAt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query segments and projects: %v", err)
	}
	defer rows.Close()

	segments := make(map[int]*entities.Segment)

	for rows.Next() {
		var segment entities.Segment
		var project entities.Project

		var projectID *int
		var projectName, projectDescription *string
		var projectCreatedAt, projectUpdatedAt, projectDeletedAt *time.Time

		err := rows.Scan(
			&segment.ID, &segment.Name, &segment.Description, &segment.CreatedAt, &segment.UpdatedAt, &segment.DeletedAt,
			&projectID, &projectName, &projectDescription, &project.Progress, &project.Url,
			&project.DateStarted, &project.DateDeadline, &projectCreatedAt, &projectUpdatedAt, &projectDeletedAt,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		if _, exists := segments[segment.ID]; !exists {
			segments[segment.ID] = &segment
		}

		if projectID != nil {
			project.ID = *projectID
			if projectName != nil {
				project.Name = *projectName
			}
			if projectDescription != nil {
				project.Description = *projectDescription
			}
			if projectCreatedAt != nil {
				project.CreatedAt = *projectCreatedAt
			}
			if projectUpdatedAt != nil {
				project.UpdatedAt = *projectUpdatedAt
			}
			if projectDeletedAt != nil {
				project.DeletedAt = projectDeletedAt
			}
			segments[segment.ID].Projects = append(segments[segment.ID].Projects, project)
		}
	}

	var result []entities.Segment
	for _, segment := range segments {
		result = append(result, *segment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return result, nil
}

func (s *Store) GetSegment(id int) (*entities.Segment, error) {
	rows, err := s.db.Query(`
		SELECT 
			seg.id AS segment_id, 
			seg.name AS segment_name, 
			seg.description AS segment_description,
			seg.createdAt AS segment_createdAt,
			seg.updatedAt AS segment_updatedAt,
			seg.deletedAt AS segment_deletedAt,
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

	segment := &entities.Segment{}
	for rows.Next() {
		var project entities.Project
		var projectID sql.NullInt64

		err := rows.Scan(
			&segment.ID, &segment.Name, &segment.Description, &segment.CreatedAt, &segment.UpdatedAt, &segment.DeletedAt,
			&projectID, &project.Name, &project.Description, &project.Progress, &project.Url,
			&project.DateStarted, &project.DateDeadline, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		if projectID.Valid {
			project.ID = int(projectID.Int64) // Convert int64 to int
			segment.Projects = append(segment.Projects, project)
		}
	}

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

	var segmentID int
	err = tx.QueryRow(
		"INSERT INTO segments (name, description) VALUES ($1, $2) RETURNING id",
		payload.Name,
		payload.Description,
	).Scan(&segmentID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return nil, err
	}

	if len(*payload.ProjectIDs) > 0 {
		for _, projectID := range *payload.ProjectIDs {
			_, err = tx.Exec(
				"INSERT INTO segments_projects (segmentId, projectId) VALUES ($1, $2)",
				segmentID,
				projectID,
			)
			if err != nil {
				// Rollback the transaction if association fails
				if rbErr := tx.Rollback(); rbErr != nil {
					return nil, fmt.Errorf("association error: %v, rollback error: %v", err, rbErr)
				}
				return nil, fmt.Errorf("failed to associate segment %d with project %d: %v", segmentID, projectID, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &entities.Segment{
		ID:          segmentID,
		Name:        payload.Name,
		Description: payload.Description,
	}, nil
}

func (s *Store) UpdateSegment(payload entities.SegmentPayload) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE segments SET name = $1, description = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3", payload.Name, payload.Description, payload.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("update error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if len(*payload.ProjectIDs) > 0 {
		_, err = tx.Exec(`
			DELETE FROM segments_projects WHERE segmentId = $1
		`, payload.ID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("delete error: %v, rollback error: %v", err, rbErr)
			}
			return fmt.Errorf("failed to delete old project associations for segment %d: %v", payload.ID, err)
		}

		for _, projectID := range *payload.ProjectIDs {
			_, err = tx.Exec(`
				INSERT INTO segments_projects (segmentId, projectId)
				VALUES ($1, $2)
			`, payload.ID, projectID)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					return fmt.Errorf("insert association error: %v, rollback error: %v", err, rbErr)
				}
				return fmt.Errorf("failed to associate segment %d with project %d: %v", payload.ID, projectID, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteSegment(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Step 1: Mark the segment as deleted
	_, err = tx.Exec("UPDATE segments SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1", id)
	if err != nil {
		// Rollback if updating the segment fails
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("delete segment error: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("failed to mark segment as deleted: %v", err)
	}

	// Step 2: Mark associated project relationships as deleted in 'segments_projects'
	_, err = tx.Exec(`
		UPDATE segments_projects SET deletedAt = CURRENT_TIMESTAMP WHERE segmentId = $1
	`, id)
	if err != nil {
		// Rollback if updating the associations fails
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("delete associations error: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("failed to mark segment-project associations as deleted: %v", err)
	}

	// Step 3: Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %v", err)
	}

	return nil
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
