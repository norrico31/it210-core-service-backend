// cmd/seed/seeders/segment_seeder.go
package seeders

import (
	"database/sql"
	"log"
	"sync"
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

	var wg sync.WaitGroup
	for _, segment := range segments {
		wg.Add(1)

		go func(segment entities.Segment) {
			defer wg.Done()

			_, err := db.Exec(`
				INSERT INTO segments (name, description, createdAt, updatedAt)
				VALUES ($1, $2, $3, $4)
			`, segment.Name, segment.Description, time.Now(), time.Now())

			if err != nil {
				log.Printf("Failed to insert segment %s: %v\n", segment.Name, err)
				return
			}
			log.Printf("Successfully inserted segment %s\n", segment.Name)
		}(segment)
	}

	return nil
}
