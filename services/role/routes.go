package role

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	router.HandleFunc("/roles", h.handleGetRoles).Methods("GET")
	router.HandleFunc("/roles", h.handleCreateRole).Methods("POST")
	router.HandleFunc("/roles/{roleId}", h.handleGetRole).Methods("GET")
	router.HandleFunc("/roles/{roleId}", h.handleUpdateRole).Methods("PUT")
	router.HandleFunc("/roles/{roleId}/restore", h.handleRestoreRole).Methods("PUT")
	router.HandleFunc("/roles/{roleId}", h.handleDeleteRole).Methods("DELETE")
}
