package tasks

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, h *Handler) {
	router.HandleFunc("/tasks", h.handleGetTasks).Methods("GET")
}
