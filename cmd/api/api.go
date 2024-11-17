package api

import (
	"database/sql"
	"log"
	"net/http"

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

// TODO: STILL NOT WORKING IN CONTAINER (PORT VARIES EVERYTIME)
func (s *APIServer) enforceGatewayOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Construct allowed host based on the config
		// allowedHost := fmt.Sprintf("%s:%s", s.config.PublicHost, s.config.GatewayPort)

		// fmt.Printf("allowedHost: %v", allowedHost)
		// fmt.Printf("rHost: %v", r.Host)

		// if r.Host == allowedHost {
		//			// Allow requests that come from the gateway (127.0.0.1:8080)
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// If the request is directly to the auth service (127.0.0.1:8081), return NOT FOUND
		// if r.Host == fmt.Sprintf("127.0.0.1:%s", s.config.GatewayPort) {
		// 	http.Error(w, "NOT FOUND", http.StatusNotFound)
		// 	return
		// }

		// Optional: Check the referer header as additional verification
		// if !strings.HasPrefix(r.Referer(), fmt.Sprintf("http://%s", allowedHost)) {
		// 	http.Error(w, "NOT FOUND", http.StatusNotFound)
		// 	return
		// }

		// Allow request to proceed if it's from the correct gateway
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
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

	log.Println("Core Service: Running on port ", s.addr)
	return http.ListenAndServe(s.addr, router)
}
