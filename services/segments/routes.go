package segments

import (
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/utils"
)

func RegisterRoutes(router *mux.Router, h *Handler) {
	utils.SecureRoute(router, "/segments", h.handleGetSegments, "GET")
	utils.SecureRoute(router, "/segments", h.handleCreateSegment, "POST")
	utils.SecureRoute(router, "/segments/{segmentId}", h.handleGetSegment, "GET")
	utils.SecureRoute(router, "/segments/{segmentId}", h.handleUpdateSegment, "PUT")
	utils.SecureRoute(router, "/segments/{segmentId}/restore", h.handleRestoreSegment, "PUT")
	utils.SecureRoute(router, "/segments/{segmentId}", h.handleDeleteSegment, "DELETE")

}
