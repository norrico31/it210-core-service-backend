package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	// router := mux.NewRouter()
	// subrouterv1 := router.PathPrefix("/api/v1").Subrouter()

	router := http.NewServeMux()
	godotenv.Load()

	// PUBLIC ROUTES
	router.HandleFunc("/api/v1/core/helloworld", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		str := fmt.Sprintf("Core Service: Hello World")
		json.NewEncoder(w).Encode(str)
	})
	log.Println("Core Service: Running on port", s.addr)
	return http.ListenAndServe(s.addr, router)
}
