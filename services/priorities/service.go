package priorities

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
	store entities.PriorityStore
}

func NewHandler(store entities.PriorityStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleGetPriorities(w http.ResponseWriter, r *http.Request) {
	prioritys, err := h.store.GetPriorities()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	w.Header().Set("Content-Type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": prioritys})
}

func (h *Handler) handleGetPriority(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["priorityId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing priority ID"))
		return
	}
	priorityId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid priority ID"))
		return
	}

	priority, err := h.store.GetPriority(priorityId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": priority})
}

func (h *Handler) handleCreatePriority(w http.ResponseWriter, r *http.Request) {
	payload := entities.PriorityPayload{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	priority, err := h.store.CreatePriority(entities.PriorityPayload{
		Name:        payload.Name,
		Description: payload.Description,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": priority})
}

func (h *Handler) handleUpdatePriority(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["priorityId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing priority ID"))
		return
	}

	priorityId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid priority ID"))
		return
	}
	var payload entities.PriorityPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	priority, err := h.store.GetPriority(priorityId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if payload.Name != "" {
		// if len(payload.Name) < 3 || len(payload.Name) > 50 {
		// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("name must be between 3 and 50 characters"))
		// 	return
		// }
		priority.Name = payload.Name
	}

	if payload.Description != "" {
		priority.Description = payload.Description
	}

	err = h.store.UpdatePriority(entities.PriorityPayload{
		ID:          priority.ID,
		Name:        priority.Name,
		Description: priority.Description,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Update Priority Successfully!"})
}

func (h *Handler) handleDeletePriority(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["priorityId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing priority ID"))
		return
	}

	priorityId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	existingPriority, err := h.store.GetPriority(priorityId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.DeletePriority(existingPriority.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]interface{}{"msg": "Delete Priority Successfully!"})
}

func (h *Handler) handleRestorePriority(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["priorityId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing priority ID"))
		return
	}

	priorityId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	existingPriority, err := h.store.GetPriority(priorityId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.RestorePriority(existingPriority.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
