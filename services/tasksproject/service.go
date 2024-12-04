package tasksproject

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
	store entities.TasksProjectStore
}

func NewHandler(store entities.TasksProjectStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleGetTasksProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["projectId"]
	if !ok {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid task id"))
		return
	}

	projectId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	tasksProject, err := h.store.GetTasksProject(projectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": tasksProject})
}

func (h *Handler) handleGetTaskProject(w http.ResponseWriter, r *http.Request) {
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
	tasksProject, err := h.store.GetTaskProject(taskId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": tasksProject})
}

func (h *Handler) handleTasksProjectCreate(w http.ResponseWriter, r *http.Request) {
	payload := entities.TasksProjectCreatePayload{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	if len(payload.Name) <= 2 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("title must be atleast 3 characters"))
		return
	}

	if payload.PriorityID == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing priority ID"))
		return
	}

	if payload.UserID == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	if payload.ProjectID == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing project ID"))
		return
	}

	task, err := h.store.TasksProjectCreate(payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": task})
}

func (h *Handler) handleTasksProjectUpdate(w http.ResponseWriter, r *http.Request) {
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

	payload := entities.TasksProjectUpdatePayload{}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	existTask, err := h.store.GetTaskProject(taskId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	if len(payload.Name) == 0 {
		payload.Name = existTask.Name
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
	if payload.ProjectID == 0 {
		payload.ProjectID = existTask.ProjectID
	}
	payload.ID = existTask.ID

	err = h.store.TasksProjectUpdate(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Update TasksProject Successfully"})
}

func (h *Handler) handleGetDeletedTasksProject(w http.ResponseWriter, r *http.Request) {
	// tasksProject, err := h.store.GetTasksProject()
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err)
	// 	return
	// }

	// w.Header().Set("Content-type", "application/json")

	// utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": tasksProject})

}

// func (h *Handler) handleGetTasksProject(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	str, ok := vars["taskId"]
// 	if !ok {
// 		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid task id"))
// 		return
// 	}

// 	taskId, err := strconv.Atoi(str)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	task, err := h.store.GetTasksProjectByID(taskId)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusNotFound, err)
// 		return
// 	}
// 	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": task})

// }

func (h *Handler) handleTasksProjectDelete(w http.ResponseWriter, r *http.Request) {
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

	err = h.store.TasksProjectDelete(taskId)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Delete TasksProject Successfully"})

}

func (h *Handler) handleTasksProjectRestore(w http.ResponseWriter, r *http.Request) {
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

	task, err := h.store.TasksProjectRestore(taskId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Restore TasksProject Successfully!", "data": task})

}

// func (h *Handler) handleTasksProjectDragNDrop(w http.ResponseWriter, r *http.Request) {
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
// 	var payload entities.TasksProjectDragNDrop

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
// 	err = h.store.TasksProjectDragNDrop(workspaceId, payload.SourceIndex, payload.DestinationIndex)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
// 		"message": "TasksProject order updated successfully",
// 	})
// }
