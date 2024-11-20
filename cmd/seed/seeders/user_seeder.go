package seeders

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/norrico31/it210-core-service-backend/entities"
	"golang.org/x/crypto/bcrypt"
)

func SeedUsers(db *sql.DB) error {
	hashPassword := func(password string) string {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}
		return string(hashedPassword)
	}

	roles := map[string]int{}
	rows, err := db.Query("SELECT id, name FROM roles")
	if err != nil {
		log.Fatalf("roles not yet seed: %v", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var roleId int
		var roleName string
		err := rows.Scan(&roleId, &roleName)
		if err != nil {
			log.Fatalf("roles not yet seed: %v", err)
			return err
		}
		roles[roleName] = roleId
	}

	admin := roles["Admin"]
	employee := roles["Employee"]

	now := time.Now()
	users := []entities.User{
		{
			FirstName: "Mary Grace",
			LastName:  "Bitmal",
			Age:       20,
			Email:     "mvbitmal@up.edu.ph",
			RoleId:    &admin,
			Password:  hashPassword("secret.123"),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			FirstName: "Chester",
			LastName:  "Francisco",
			Age:       19,
			RoleId:    &admin,
			Email:     "cgfrancisco@up.edu.ph",
			Password:  hashPassword("secret.123"),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			FirstName: "Norrico Gerald",
			LastName:  "Biason",
			Age:       18,
			RoleId:    &employee,
			Email:     "nmbiason@up.edu.ph",
			Password:  hashPassword("secret.123"),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(users))

	for _, user := range users {
		go func(u entities.User) {
			defer wg.Done()

			_, err := db.Exec(`
			INSERT INTO users (firstName, lastName, age, email, roleId, password, createdAt, updatedAt) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				u.FirstName,
				u.LastName,
				u.Age,
				u.Email,
				u.RoleId,
				u.Password,
				u.CreatedAt,
				u.UpdatedAt)
			if err != nil {
				log.Printf("Failed to insert user %s: %v", u.Email, err)
			} else {
				log.Printf("Inserted user: %s", u.Email)
			}
		}(user)
	}
	wg.Wait()
	log.Println("Users table seeded successfully.")
	return nil
}
