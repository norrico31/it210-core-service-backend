// cmd/seed/seeders/segment_seeder.go
package seeders

import (
	"database/sql"
	"log"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
)

func SeedSegments(db *sql.DB) error {
	segments := []entities.Segment{
		{
			Name:        "SIYENSIKAT",
			Description: "siyensikat",
		},
		{
			Name:        "EXPERTALK",
			Description: "expertalk",
		},
		{
			Name:        "NEGOSIYENSIYA",
			Description: "negosiyensiya",
		},
		{
			Name:        "BANTAY BULKAN",
			Description: "bantay bulkan",
		},
		{
			Name:        "RAPIDDOST",
			Description: "rapiddost",
		},
	}

	for _, segment := range segments {
		_, err := db.Exec(`
				INSERT INTO segments (name, description, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4)
			`, segment.Name, segment.Description, time.Now(), time.Now())

		if err != nil {
			log.Printf("Failed to insert segment %s: %v\n", segment.Name, err)
			return err
		}
		log.Printf("Successfully inserted segment %s\n", segment.Name)
	}

	return nil
}
