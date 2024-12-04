package workspaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/norrico31/it210-core-service-backend/entities"
	"github.com/norrico31/it210-core-service-backend/utils"
)

type Handler struct {
	store entities.WorkspaceStore
}

func NewHandler(store entities.WorkspaceStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleGetWorkspaces(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// str, ok := vars["projectId"]
	// if !ok {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing workspace ID"))
	// 	return
	// }
	// projectId, err := strconv.Atoi(str)
	workspaces, err := h.store.GetWorkspaces()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": workspaces})
}

func (h *Handler) handleGetWorkspace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["projectId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing project ID"))
		return
	}
	projectId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	workspace, err := h.store.GetWorkspace(projectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": workspace})
}

func (h *Handler) handleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
	payload := entities.WorkspacePayload{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	workspace, err := h.store.CreateWorkspace(payload)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": workspace})
}

func (h *Handler) handleUpdateWorkspace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["workspaceId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing workspace ID"))
		return
	}

	workspaceId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid workspace ID"))
		return
	}
	var payload entities.WorkspacePayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	payload.ID = workspaceId
	err = h.store.UpdateWorkspace(payload)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Update Workspace Successfully!"})
}

// This is an example route handler for the TaskDragNDrop function
func (h *Handler) handleTaskDragNDrop(w http.ResponseWriter, r *http.Request) {
	// Extract the workspaceId from the URL parameter
	vars := mux.Vars(r)
	str, ok := vars["workspaceId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing workspace ID"))
		return
	}

	workspaceId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid workspace ID"))
		return
	}

	// Parse the request body to get sourceIndex and destinationIndex
	var payload entities.TaskDragNDrop

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	print("source ", payload.SourceIndex)
	print("destination ", payload.DestinationIndex)
	print("workspaceId ", workspaceId)

	// Call the TaskDragNDrop method
	err = h.store.TaskDragNDrop(workspaceId, payload.SourceIndex, payload.DestinationIndex)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update task order: %v", err), http.StatusInternalServerError)
		return
	}

	// Return a successful response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Task order updated successfully")
}

func (h *Handler) handleDeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// str, ok := vars["workspaceId"]
	// if !ok {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing workspace ID"))
	// 	return
	// }

	// workspaceId, err := strconv.Atoi(str)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
	// 	return
	// }

	// existingWorkspace, err := h.store.GetWorkspace(workspaceId)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err)
	// 	return
	// }

	// err = h.store.DeleteWorkspace(existingWorkspace.ID)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, err)
	// 	return
	// }

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Delete Workspace Successfully!"})
}

func (h *Handler) handleRestoreWorkspace(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// str, ok := vars["workspaceId"]
	// if !ok {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing workspace ID"))
	// 	return
	// }

	// workspaceId, err := strconv.Atoi(str)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
	// 	return
	// }

	// existingWorkspace, err := h.store.GetWorkspace(workspaceId)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err)
	// 	return
	// }

	// err = h.store.RestoreWorkspace(existingWorkspace.ID)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, err)
	// 	return
	// }

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Restore Workspace Successfully!"})
}
