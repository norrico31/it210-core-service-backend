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
