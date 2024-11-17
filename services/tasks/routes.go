package tasks

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, h *Handler) {
	router.HandleFunc("/tasks", h.handleGetTasks).Methods("GET")
	router.HandleFunc("/tasks", h.handleTaskCreate).Methods("POST")
	router.HandleFunc("/tasks/{taskId}", h.handleGetTask).Methods("GET")
	router.HandleFunc("/tasks/{taskId}", h.handleTaskUpdate).Methods("PUT")
	router.HandleFunc("/tasks/{taskId}", h.handleTaskDelete).Methods("DELETE")
	router.HandleFunc("/tasks/{taskId}/restore", h.handleTaskRestore).Methods("PUT")
}
