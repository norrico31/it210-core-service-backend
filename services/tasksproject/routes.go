package tasksproject

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	utils.SecureRoute(router, "/tasks/{projectId}", h.handleGetTasksProject, "GET")
	utils.SecureRoute(router, "/tasks/{taskId}/task", h.handleGetTaskProject, "GET")
	utils.SecureRoute(router, "/tasks/{projectId}", h.handleTasksProjectCreate, "POST")
	utils.SecureRoute(router, "/tasks/deleted", h.handleGetDeletedTasksProject, "GET")
	// utils.SecureRoute(router, "/tasks/reorder", h.handleTaskDragNDrop, "PUT")
	utils.SecureRoute(router, "/tasks/{taskId}", h.handleTasksProjectUpdate, "PUT")
	utils.SecureRoute(router, "/tasks/{taskId}", h.handleTasksProjectDelete, "DELETE")
	utils.SecureRoute(router, "/tasks/{taskId}/restore", h.handleTasksProjectRestore, "PUT")
}
