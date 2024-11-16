package tasks

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, h *Handler) {
	router.HandleFunc("/tasks", h.handleGetTasks).Methods("GET")
	router.HandleFunc("/tasks", h.handleTaskCreate).Methods("POST")
	router.HandleFunc("/tasks/{taskId}", h.handleGetTask).Methods("GET")
}
