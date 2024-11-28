package users

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	utils.SecureRoute(router, "/users", h.handleGetUsers, "GET")
	utils.SecureRoute(router, "/users", h.handleCreateUser, "POST")
	utils.SecureRoute(router, "/users/{userId}", h.handleGetUser, "GET")
	utils.SecureRoute(router, "/users/{userId}", h.HandleUpdateUser, "PUT")
	utils.SecureRoute(router, "/users/{userId}", h.HandleDeleteUser, "DELETE")
	utils.SecureRoute(router, "/users/{userId}/restore", h.handleRestoreUser, "PUT")
	utils.SecureRoute(router, "/users/logout/{userId}", h.handleLogout, "POST")
	// utils.SecureRoute(router, "/user/create", h.handleLogout, "POST")
	// router.HandleFunc("/users/register", h.handleCreateUser).Methods("POST")
}
