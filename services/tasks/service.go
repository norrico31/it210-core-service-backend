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

	payload := entities.TaskUpdatePayload{}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	existTask, err := h.store.GetTask(taskId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	if len(payload.Title) == 0 {
		payload.Title = existTask.Title
	}

	if payload.Description == "" {
		payload.Description = existTask.Description
	}
	if payload.UserID == 0 {
		payload.UserID = *existTask.UserID
	}

	if payload.PriorityID == 0 {
		payload.PriorityID = existTask.PriorityID
	}

	if payload.WorkspaceID == 0 {
		payload.WorkspaceID = existTask.WorkspaceID
	}
	payload.ID = existTask.ID

	err = h.store.TaskUpdate(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Update Task Successfully"})

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

	err = h.store.TaskDelete(taskId)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Delete Task Successfully"})

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

// func (h *Handler) handleTaskDragNDrop(w http.ResponseWriter, r *http.Request) {
// 	// Parse query parameter for workspaceId
// 	workspaceIdStr := r.URL.Query().Get("workspaceId")
// 	if workspaceIdStr == "" {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing workspaceId query parameter"))
// 		return
// 	}
// 	print("pumapasok ba siya dito?")
// 	workspaceId, err := strconv.Atoi(workspaceIdStr)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid workspaceId"))
// 		return
// 	}

// 	// Parse JSON payload
// 	var payload entities.TaskDragNDrop

// 	if err := utils.ParseJSON(r, &payload); err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON payload"))
// 		return
// 	}

// 	// Validate payload
// 	if payload.SourceIndex < 0 || payload.DestinationIndex < 0 {
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("sourceIndex and destinationIndex must be non-negative"))
// 		return
// 	}

// 	// Call store method to update task order
// 	err = h.store.TaskDragNDrop(workspaceId, payload.SourceIndex, payload.DestinationIndex)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
// 		"message": "Task order updated successfully",
// 	})
// }
