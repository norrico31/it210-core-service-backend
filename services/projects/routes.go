package projects

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, h *Handler) {
	router.HandleFunc("/projects", h.handleGetProjects).Methods("GET")
}
