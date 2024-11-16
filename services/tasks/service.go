package tasks

import (
	"net/http"

	"github.com/norrico31/it210-core-service-backend/entities"
	"github.com/norrico31/it210-core-service-backend/utils"
)

type Handler struct {
	store entities.TaskStore
}

func NewHandler(store entities.TaskStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {

	roles, err := h.store.GetTasks()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": roles})
}
