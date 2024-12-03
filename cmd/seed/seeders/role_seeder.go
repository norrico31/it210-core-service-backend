// cmd/seed/seeders/role_seeder.go
package seeders

import (
	"database/sql"
	"log"
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

	for _, role := range roles {
		_, err := db.Exec(`
				INSERT INTO roles (name, description, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4)
			`, role.Name, role.Description, time.Now(), time.Now())

		if err != nil {
			log.Printf("Failed to insert role %s: %v\n", role.Name, err)
			return err
		}
		log.Printf("Successfully inserted role %s\n", role.Name)
	}

	return nil
}
