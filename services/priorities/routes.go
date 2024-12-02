package priorities

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	utils.SecureRoute(router, "/priorities", h.handleGetPriorities, "GET")
	utils.SecureRoute(router, "/priorities", h.handleCreatePriority, "POST")
	utils.SecureRoute(router, "/priorities/{priorityId}", h.handleGetPriority, "GET")
	utils.SecureRoute(router, "/priorities/{priorityId}", h.handleUpdatePriority, "PUT")
	utils.SecureRoute(router, "/priorities/{priorityId}/restore", h.handleRestorePriority, "PUT")
	utils.SecureRoute(router, "/priorities/{priorityId}", h.handleDeletePriority, "DELETE")

}
