package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/config"
	"github.com/norrico31/it210-core-service-backend/services/projects"
	"github.com/norrico31/it210-core-service-backend/services/roles"
	"github.com/norrico31/it210-core-service-backend/services/tasks"
	"github.com/norrico31/it210-core-service-backend/services/users"
)

type APIServer struct {
	addr   string
	db     *sql.DB
	config config.Config
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

// Request logging middleware
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details: method, URL, and client IP
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// TODO: STILL NOT WORKING IN CONTAINER (PORT VARIES EVERYTIME)
func (s *APIServer) enforceGatewayOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Your existing logic for enforcing gateway origin
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// Apply the request logging middleware
	router.Use(logRequest)

	// Apply the gateway origin enforcement middleware
	router.Use(s.enforceGatewayOrigin)

	subrouterv1 := router.PathPrefix("/api/v1/core").Subrouter()

	roleStore := roles.NewStore(s.db)
	roleHandler := roles.NewHandler(roleStore)
	roles.RegisterRoutes(subrouterv1, roleHandler)

	usersStore := users.NewStore(s.db)
	usersHandler := users.NewHandler(usersStore)
	users.RegisterRoutes(subrouterv1, usersHandler)

	taskStore := tasks.NewStore(s.db)
	taskHandler := tasks.NewHandler(taskStore)
	tasks.RegisterRoutes(subrouterv1, taskHandler)

	projectStore := projects.NewStore(s.db)
	projecthandler := projects.NewHandler(projectStore)
	projects.RegisterRoutes(subrouterv1, projecthandler)

	// CORS configuration
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // You can replace "*" with specific allowed origins if needed
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router)

	// Create and start the server
	server := &http.Server{
		Addr:           ":8080",
		Handler:        corsHandler,
		MaxHeaderBytes: 1 << 20, // 1 MB for header size, adjust as needed
	}
	log.Println("Core Service: Running on port ", s.addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
