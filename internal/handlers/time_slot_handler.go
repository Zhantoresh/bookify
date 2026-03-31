package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bookify/internal/service"
)

type TimeSlotHandler struct {
	timeSlotService *service.TimeSlotService
}

func NewTimeSlotHandler(timeSlotService *service.TimeSlotService) *TimeSlotHandler {
	return &TimeSlotHandler{
		timeSlotService: timeSlotService,
	}
}

func (h *TimeSlotHandler) HandleTimeSlots(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTimeSlot(w, r)
	case http.MethodGet:
		h.getMyTimeSlots(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TimeSlotHandler) HandleTimeSlotsWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.updateTimeSlot(w, r)
	case http.MethodDelete:
		h.deleteTimeSlot(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TimeSlotHandler) createTimeSlot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		respondError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	var input struct {
		Time string `json:"time"` // Format: 2006-01-02T15:04:05Z
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid input")
		return
	}

	slotTime, err := time.Parse(time.RFC3339, input.Time)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid time format, use RFC3339")
		return
	}

	slot, err := h.timeSlotService.CreateTimeSlot(userID, slotTime)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create time slot")
		return
	}

	respondJSON(w, http.StatusCreated, slot)
}

func (h *TimeSlotHandler) getMyTimeSlots(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		respondError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	slots, err := h.timeSlotService.GetMyTimeSlots(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch time slots")
		return
	}

	respondJSON(w, http.StatusOK, slots)
}

func (h *TimeSlotHandler) updateTimeSlot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		respondError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/time-slots/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid time slot id")
		return
	}

	var input struct {
		Time string `json:"time"` // Format: 2006-01-02T15:04:05Z
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid input")
		return
	}

	slotTime, err := time.Parse(time.RFC3339, input.Time)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid time format, use RFC3339")
		return
	}

	err = h.timeSlotService.UpdateTimeSlot(id, userID, slotTime)
	if err != nil {
		if errors.Is(err, service.ErrTimeSlotNotFound) {
			respondError(w, http.StatusNotFound, "time slot not found")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			respondError(w, http.StatusForbidden, "you can only manage your own time slots")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update time slot")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "time slot updated"})
}

func (h *TimeSlotHandler) deleteTimeSlot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		respondError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/time-slots/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid time slot id")
		return
	}

	err = h.timeSlotService.DeleteTimeSlot(id, userID)
	if err != nil {
		if errors.Is(err, service.ErrTimeSlotNotFound) {
			respondError(w, http.StatusNotFound, "time slot not found")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			respondError(w, http.StatusForbidden, "you can only manage your own time slots")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete time slot")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "time slot deleted"})
}
