package projects

import (
	"net/http"

	"github.com/norrico31/it210-core-service-backend/entities"
	"github.com/norrico31/it210-core-service-backend/utils"
)

type Handler struct {
	store entities.ProjectStore
}

func NewHandler(store entities.ProjectStore) *Handler {
	return &Handler{store: store}
}

func (s *Handler) handleGetProjects(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"hi": "HELLO WORLD FROM GETPROJECTS"})
}
