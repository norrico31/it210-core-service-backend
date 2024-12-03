package workspaces

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	utils.SecureRoute(router, "/workspaces", h.handleGetWorkspaces, "GET")
	utils.SecureRoute(router, "/workspaces/{projectId}", h.handleGetWorkspaces, "GET")
	// utils.SecureRoute(router, "/workspaces", h.handleCreateWorkspace, "POST")
	// utils.SecureRoute(router, "/workspaces/{workspaceId}", h.handleGetWorkspace, "GET")
	utils.SecureRoute(router, "/workspaces/{workspaceId}", h.handleUpdateWorkspace, "PUT")
	utils.SecureRoute(router, "/workspaces/{workspaceId}/restore", h.handleRestoreWorkspace, "PUT")
	utils.SecureRoute(router, "/workspaces/{workspaceId}", h.handleDeleteWorkspace, "DELETE")

}
