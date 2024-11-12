package role

import (
	"net/http"

	_ "github.com/go-playground/validator/v10"
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
