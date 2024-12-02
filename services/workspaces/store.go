package workspaces

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

func (s *Store) GetWorkspaces() ([]entities.Workspace, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, createdAt, updatedAt
			FROM workspaces
		WHERE deletedAt IS NULL
		ORDER BY createdAt DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspaces: %v", err)
	}
	defer rows.Close()

	workspaces := []entities.Workspace{}

	for rows.Next() {
		workspace := entities.Workspace{}

		err := scanRowIntoWorkspace(rows, &workspace)
		if err != nil {
			log.Printf("Failed to scan workspace: %v", err)
			continue
		}
		workspaces = append(workspaces, workspace)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over workspace rows: %v", err)
	}
	return workspaces, nil
}

func (s *Store) GetWorkspace(id int) (*entities.Workspace, error) {
	workspace := entities.Workspace{}
	err := s.db.QueryRow("SELECT id, name, description, createdAt, updatedAt FROM workspaces WHERE deletedAt IS NULL AND id = $1", id).Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Description,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("workspace not found")
	}

	if workspace.ID == 0 {
		return nil, fmt.Errorf("workspace not found")
	}
	return &workspace, nil
}

func (s *Store) CreateWorkspace(payload entities.WorkspacePayload) (*entities.Workspace, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO workspaces (name, description) VALUES ($1, $2)", payload.Name, payload.Description)
	if err != nil {
		// If there's an error, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
		}
		return nil, err
	}

	// Commit the transaction if all went well
	if err = tx.Commit(); err != nil {
		return &entities.Workspace{Name: payload.Name, Description: payload.Description}, err
	}

	return &entities.Workspace{Name: payload.Name, Description: payload.Description}, err
}

func (s *Store) UpdateWorkspace(payload entities.WorkspacePayload) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE workspaces SET name = $1, description = $2, updatedAt = CURRENT_TIMESTAMP WHERE id = $3", payload.Name, payload.Description, payload.ID)
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

func (s *Store) DeleteWorkspace(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE workspaces SET deletedAt = CURRENT_TIMESTAMP WHERE id = $1", id)

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

func (s *Store) RestoreWorkspace(id int) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE workspaces SET deletedAt = NULL WHERE id = $1", id)
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

func scanRowIntoWorkspace(rows *sql.Rows, workspace *entities.Workspace) error {
	return rows.Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Description,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
}
