package tasksproject

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	utils.SecureRoute(router, "/tasks", h.handleTasksProjectCreate, "POST")
	utils.SecureRoute(router, "/tasks/{taskId}/delete", h.handleGetDeletedTasksProject, "DELETE")
	utils.SecureRoute(router, "/tasks/{projectId}", h.handleGetTasksProject, "GET")
	utils.SecureRoute(router, "/tasks/{taskId}/task", h.handleGetTaskProject, "GET")
	utils.SecureRoute(router, "/tasks/{taskId}/task", h.handleTasksProjectUpdate, "PUT")
	utils.SecureRoute(router, "/tasks/{taskId}", h.handleTasksProjectUpdate, "PUT")
	utils.SecureRoute(router, "/tasks/{taskId}", h.handleTasksProjectDelete, "DELETE")
	utils.SecureRoute(router, "/tasks/{taskId}/restore", h.handleTasksProjectRestore, "PUT")
}
