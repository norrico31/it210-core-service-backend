package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
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

func (h *Handler) handleGetDeletedTasks(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	payload := entities.TaskCreatePayload{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	if len(payload.Title) <= 2 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("title must be atleast 3 characters"))
		return
	}

	if payload.PriorityID == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing priority ID"))
		return
	}

	if payload.WorkspaceID == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing workspace ID"))
		return
	}

	task, err := h.store.TaskCreate(payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": task})

}

func (h *Handler) handleTaskUpdate(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleTaskDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["taskId"]
	if !ok {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid task ID"))
		return
	}

	taskId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	existTask, err := h.store.GetTask(taskId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	task, err := h.store.TaskDelete(existTask.ID)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Delete Task Successfully", "data": task})

}

func (h *Handler) handleTaskRestore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["taskId"]
	if !ok {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid task ID"))
		return
	}

	taskId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	task, err := h.store.TaskRestore(taskId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Restore Task Successfully!", "data": task})

}
