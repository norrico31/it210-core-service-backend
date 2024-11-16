package projects

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, h *Handler) {
	router.HandleFunc("/projects", h.handleGetProjects).Methods("GET")
	router.HandleFunc("/projects", h.handleProjectCreate).Methods("POST")
	router.HandleFunc("/projects/{projectId}", h.handleGetProject).Methods("GET")
	router.HandleFunc("/projects/{projectId}", h.handleProjectUpdate).Methods("PUT")
}
