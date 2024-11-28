package tasks

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	utils.SecureRoute(router, "/tasks", h.handleGetTasks, "GET")
	utils.SecureRoute(router, "/tasks/deleted", h.handleGetDeletedTasks, "GET")
	utils.SecureRoute(router, "/tasks", h.handleTaskCreate, "POST")
	utils.SecureRoute(router, "/tasks/{taskId}", h.handleGetTask, "GET")
	utils.SecureRoute(router, "/tasks/{taskId}", h.handleTaskUpdate, "PUT")
	utils.SecureRoute(router, "/tasks/{taskId}", h.handleTaskDelete, "DELETE")
	utils.SecureRoute(router, "/tasks/{taskId}/restore", h.handleTaskRestore, "PUT")

}
