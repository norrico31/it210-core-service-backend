package users

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	router.HandleFunc("/users", h.handleGetUsers).Methods("GET")
	router.HandleFunc("/users/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/users/logout/{userId}", h.handleLogout).Methods("POST")
	router.HandleFunc("/users/{userId}", h.handleGetUser).Methods("GET")
	router.HandleFunc("/users/{userId}", h.HandleUpdateUser).Methods("PUT")
	router.HandleFunc("/users/{userId}", h.HandleDeleteUser).Methods("DELETE")
}
