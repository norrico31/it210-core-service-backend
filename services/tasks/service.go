package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	tasks, err := h.store.GetTasks()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": tasks})
}

func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["taskId"]
	if !ok {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid task id"))
		return
	}

	taskId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	task, err := h.store.GetTask(taskId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": task})

}
