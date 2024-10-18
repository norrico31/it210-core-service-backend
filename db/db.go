package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/norrico31/it210-core-service-backend/config"
)

func NewPostgresStorage() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%d sslmode=disable",
		config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBAddress, config.Envs.DBName, 5432)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
