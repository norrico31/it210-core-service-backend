package projects

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
	store entities.ProjectStore
}

func NewHandler(store entities.ProjectStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleGetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.GetProjects("IS NULL")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": projects})
}

func (h *Handler) handleGetProjectDeleted(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.GetProjects("IS NOT NULL")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": projects})
}

func (h *Handler) handleGetProject(w http.ResponseWriter, r *http.Request) {
	for header := range r.Header {
		if header == "X-User-Id" {
			userId := r.Header.Get("X-User-Id")
			fmt.Printf("userid: %s \n", userId)
		}
	}
	vars := mux.Vars(r)
	str, ok := vars["projectId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	projectId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid project Id"))
		return
	}
	proj, err := h.store.GetProject(projectId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, proj)
}

func (h *Handler) handleProjectCreate(w http.ResponseWriter, r *http.Request) {
	payload := entities.ProjectCreatePayload{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	proj, err := h.store.ProjectCreate(payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, proj)
}

// TODO
func (h *Handler) handleProjectUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["projectId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	projectId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid project Id"))
		return
	}

	payload := entities.ProjectUpdatePayload{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	existProj, err := h.store.GetProject(projectId)

	if payload.Name != "" {
		existProj.Name = payload.Name
	}
	if payload.Description != "" {
		existProj.Description = payload.Description
	}
	if payload.Progress != nil {
		existProj.Progress = payload.Progress
	}
	if payload.DateStarted != "" {
		dateStarted, err := normalizeDate(payload.DateStarted)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid date format for Date Started"))
			return
		}
		existProj.DateStarted = &dateStarted
	}
	if payload.DateDeadline != "" {
		dateDeadline, err := normalizeDate(payload.DateDeadline)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid date format for Date Deadline"))
			return
		}
		existProj.DateDeadline = &dateDeadline
	}

	if payload.Url != nil {
		existProj.Url = payload.Url
	}

	// Ensure non-nullable fields are not nullified
	if payload.StatusID != 0 {
		existProj.StatusID = payload.StatusID
	}
	if payload.SegmentID != 0 {
		existProj.SegmentID = payload.SegmentID
	}
	var userIDs []int
	if payload.UserIDs != nil {
		userIDs = *payload.UserIDs
	} else {
		for _, user := range existProj.Users {
			userIDs = append(userIDs, user.ID)
		}
	}

	err = h.store.ProjectUpdate(projectId, payload, userIDs)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Update Successfully!"})
}

func (h *Handler) handleProjectDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["projectId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	projectId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	// existProj, err := h.store.GetProject(projectId)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusNotFound, err)
	// 	return
	// }

	project, err := h.store.ProjectDelete(projectId)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Delete Project Successfully", "data": project})
}

func (h *Handler) handleProjectRestore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["projectId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	projectId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid project ID"))
		return
	}

	project, err := h.store.ProjectRestore(projectId)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Restore Project Successfully!", "data": project})
}
