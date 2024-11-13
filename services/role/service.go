package role

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
	store entities.RoleStore
}

func NewHandler(store entities.RoleStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleGetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.store.GetRoles()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	w.Header().Set("Content-Type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": roles})
}

func (h *Handler) handleGetRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["roleId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing role ID"))
		return
	}
	roleId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid role ID"))
		return
	}

	role, err := h.store.GetRole(roleId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, role)
}

func (h *Handler) handleCreateRole(w http.ResponseWriter, r *http.Request) {
	payload := entities.Role{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.store.CreateRole(entities.Role{
		Name:        payload.Name,
		Description: payload.Description,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUpdateRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["roleId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing role ID"))
		return
	}

	roleId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid role ID"))
		return
	}

	existingRole, err := h.store.GetRole(roleId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload = entities.Role{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	role := entities.Role{
		ID:          existingRole.ID,
		Name:        payload.Name,
		Description: payload.Description,
	}

	err = h.store.UpdateRole(role)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleDeleteRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["roleId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing role ID"))
		return
	}

	roleId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	existingRole, err := h.store.GetRole(roleId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.DeleteRole(existingRole.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleRestoreRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["roleId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing role ID"))
		return
	}

	roleId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	existingRole, err := h.store.GetRole(roleId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.RestoreRole(existingRole.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
