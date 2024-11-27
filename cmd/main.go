package main

import (
	"fmt"
	"log"

	"github.com/norrico31/it210-core-service-backend/cmd/api"
	"github.com/norrico31/it210-core-service-backend/db"
)

func main() {
	db, err := db.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("PostgreSQL connection established!")
	server := api.NewApiServer(":8080", db)
	server.Run()
}
