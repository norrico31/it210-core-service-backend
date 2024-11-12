package role

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	router.HandleFunc("/roles", h.handleGetRoles).Methods("GET")
}
