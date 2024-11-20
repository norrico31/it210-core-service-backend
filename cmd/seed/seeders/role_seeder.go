// cmd/seed/seeders/role_seeder.go
package seeders

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedRoles(db *sql.DB) error {
	roles := []entities.Role{
		{
			Name:        "Admin",
			Description: "administrator",
		},
		{
			Name:        "Employee",
			Description: "employee",
		},
		{
			Name:        "Manager",
			Description: "manager",
		},
	}

	var wg sync.WaitGroup
	for _, role := range roles {
		wg.Add(1)

		go func(role entities.Role) {
			defer wg.Done()

			_, err := db.Exec(`
				INSERT INTO roles (name, description, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4)
			`, role.Name, role.Description, time.Now(), time.Now())

			if err != nil {
				log.Printf("Failed to insert role %s: %v\n", role.Name, err)
				return
			}
			log.Printf("Successfully inserted role %s\n", role.Name)
		}(role)
	}

	return nil
}
