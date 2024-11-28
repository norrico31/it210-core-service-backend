package projects

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	utils.SecureRoute(router, "/projects", h.handleGetProjects, "GET")
	utils.SecureRoute(router, "/projects", h.handleProjectCreate, "POST")
	utils.SecureRoute(router, "/projects/deleted", h.handleGetProjectDeleted, "GET")
	utils.SecureRoute(router, "/projects/{projectId}", h.handleGetProject, "GET")
	utils.SecureRoute(router, "/projects/{projectId}", h.handleProjectUpdate, "PUT")
	utils.SecureRoute(router, "/projects/{projectId}", h.handleProjectDelete, "DELETE")
	utils.SecureRoute(router, "/projects/{projectId}/restore", h.handleProjectRestore, "PUT")
}
