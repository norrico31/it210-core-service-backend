package roles

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	utils.SecureRoute(router, "/roles", h.handleGetRoles, "GET")
	utils.SecureRoute(router, "/roles", h.handleCreateRole, "POST")
	utils.SecureRoute(router, "/roles/{roleId}", h.handleGetRole, "GET")
	utils.SecureRoute(router, "/roles/{roleId}", h.handleUpdateRole, "PUT")
	utils.SecureRoute(router, "/roles/{roleId}/restore", h.handleRestoreRole, "PUT")
	utils.SecureRoute(router, "/roles/{roleId}", h.handleDeleteRole, "DELETE")

}
