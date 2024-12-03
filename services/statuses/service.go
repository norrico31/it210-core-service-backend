package statuses

import (
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
	store entities.StatusStore
}

func NewHandler(store entities.StatusStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleGetStatuses(w http.ResponseWriter, r *http.Request) {
	statuses, err := h.store.GetStatuses()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": statuses})
}

func (h *Handler) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["statusId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing statuses ID"))
		return
	}
	statusesId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid statuses ID"))
		return
	}

	statuses, err := h.store.GetStatus(statusesId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": statuses})
}

func (h *Handler) handleCreateStatus(w http.ResponseWriter, r *http.Request) {
	payload := entities.StatusPayload{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	statuses, err := h.store.CreateStatus(entities.StatusPayload{
		Name:        payload.Name,
		Description: payload.Description,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": statuses})
}

func (h *Handler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["statusId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing statuses ID"))
		return
	}

	statusesId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid statuses ID"))
		return
	}
	var payload entities.StatusPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	statuses, err := h.store.GetStatus(statusesId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if payload.Name != "" {
		statuses.Name = payload.Name
	}

	if payload.Description != "" {
		statuses.Description = payload.Description
	}

	err = h.store.UpdateStatus(entities.StatusPayload{
		ID:          statuses.ID,
		Name:        statuses.Name,
		Description: statuses.Description,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Update Status Successfully!"})
}

func (h *Handler) handleDeleteStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["statusId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing statuses ID"))
		return
	}

	statusesId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	existingStatus, err := h.store.GetStatus(statusesId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.DeleteStatus(existingStatus.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Delete Status Successfully!"})
}

func (h *Handler) handleRestoreStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["statusesId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing statuses ID"))
		return
	}

	statusesId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	existingStatus, err := h.store.GetStatus(statusesId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.RestoreStatus(existingStatus.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
