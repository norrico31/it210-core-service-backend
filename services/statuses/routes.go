package statuses

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	utils.SecureRoute(router, "/statuses", h.handleGetStatuses, "GET")
	utils.SecureRoute(router, "/statuses", h.handleCreateStatus, "POST")
	utils.SecureRoute(router, "/statuses/{statusId}", h.handleGetStatus, "GET")
	utils.SecureRoute(router, "/statuses/{statusId}", h.handleUpdateStatus, "PUT")
	utils.SecureRoute(router, "/statuses/{statusId}/restore", h.handleRestoreStatus, "PUT")
	utils.SecureRoute(router, "/statuses/{statusId}", h.handleDeleteStatus, "DELETE")

}
