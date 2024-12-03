package segments

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
	store entities.SegmentsStore
}

func NewHandler(store entities.SegmentsStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleGetSegments(w http.ResponseWriter, r *http.Request) {
	segments, err := h.store.GetSegments()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": segments})
}

func (h *Handler) handleGetSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["segmentId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing segment ID"))
		return
	}
	segmentId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid segment ID"))
		return
	}

	segment, err := h.store.GetSegment(segmentId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": segment})
}

func (h *Handler) handleCreateSegment(w http.ResponseWriter, r *http.Request) {
	payload := entities.SegmentPayload{}

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	segment, err := h.store.CreateSegment(entities.SegmentPayload{
		Name:        payload.Name,
		Description: payload.Description,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": segment})
}

func (h *Handler) handleUpdateSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["segmentId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing segment ID"))
		return
	}

	segmentId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid segment ID"))
		return
	}
	var payload entities.SegmentPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	segment, err := h.store.GetSegment(segmentId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if payload.Name != "" {
		// if len(payload.Name) < 3 || len(payload.Name) > 50 {
		// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("name must be between 3 and 50 characters"))
		// 	return
		// }
		segment.Name = payload.Name
	}

	if payload.Description != "" {
		segment.Description = payload.Description
	}

	err = h.store.UpdateSegment(entities.SegmentPayload{
		ID:          segment.ID,
		Name:        segment.Name,
		Description: segment.Description,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"msg": "Update Segment Successfully!"})
}

func (h *Handler) handleDeleteSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["segmentId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing segment ID"))
		return
	}

	segmentId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	existingSegment, err := h.store.GetSegment(segmentId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.DeleteSegment(existingSegment.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]interface{}{"msg": "Delete Segment Successfully!"})
}

func (h *Handler) handleRestoreSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["segmentId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing segment ID"))
		return
	}

	segmentId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	existingSegment, err := h.store.GetSegment(segmentId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.RestoreSegment(existingSegment.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
