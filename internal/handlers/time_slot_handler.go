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
		logAndRespondError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed for /time-slots", nil)
	}
}

func (h *TimeSlotHandler) HandleTimeSlotsWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.updateTimeSlot(w, r)
	case http.MethodDelete:
		h.deleteTimeSlot(w, r)
	default:
		logAndRespondError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed for /time-slots/{id}", nil)
	}
}

func (h *TimeSlotHandler) createTimeSlot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		logAndRespondError(w, http.StatusUnauthorized, "user id not found in context", "user_id not found in request context", nil)
		return
	}

	var input struct {
		Time string `json:"time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logAndRespondError(w, http.StatusBadRequest, "invalid input", "failed to decode create time slot request", err)
		return
	}

	slotTime, err := time.Parse(time.RFC3339, input.Time)
	if err != nil {
		logAndRespondError(w, http.StatusBadRequest, "invalid time format, use RFC3339", "failed to parse time slot time", err)
		return
	}

	slot, err := h.timeSlotService.CreateTimeSlot(userID, slotTime)
	if err != nil {
		logAndRespondError(w, http.StatusInternalServerError, "failed to create time slot", "failed to create time slot", err)
		return
	}

	respondJSON(w, http.StatusCreated, slot)
}

func (h *TimeSlotHandler) getMyTimeSlots(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		logAndRespondError(w, http.StatusUnauthorized, "user id not found in context", "user_id not found in request context", nil)
		return
	}

	slots, err := h.timeSlotService.GetMyTimeSlots(userID)
	if err != nil {
		logAndRespondError(w, http.StatusInternalServerError, "failed to fetch time slots", "failed to fetch user time slots", err)
		return
	}

	respondJSON(w, http.StatusOK, slots)
}

func (h *TimeSlotHandler) updateTimeSlot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		logAndRespondError(w, http.StatusUnauthorized, "user id not found in context", "user_id not found in request context", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/time-slots/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logAndRespondError(w, http.StatusBadRequest, "invalid time slot id", "failed to parse time slot id", err)
		return
	}

	var input struct {
		Time string `json:"time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logAndRespondError(w, http.StatusBadRequest, "invalid input", "failed to decode update time slot request", err)
		return
	}

	slotTime, err := time.Parse(time.RFC3339, input.Time)
	if err != nil {
		logAndRespondError(w, http.StatusBadRequest, "invalid time format, use RFC3339", "failed to parse updated time slot time", err)
		return
	}

	err = h.timeSlotService.UpdateTimeSlot(id, userID, slotTime)
	if err != nil {
		if errors.Is(err, service.ErrTimeSlotNotFound) {
			logAndRespondError(w, http.StatusNotFound, "time slot not found", "time slot not found during update", err)
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			logAndRespondError(w, http.StatusForbidden, "you can only manage your own time slots", "forbidden time slot update attempt", err)
			return
		}
		logAndRespondError(w, http.StatusInternalServerError, "failed to update time slot", "failed to update time slot", err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "time slot updated"})
}

func (h *TimeSlotHandler) deleteTimeSlot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		logAndRespondError(w, http.StatusUnauthorized, "user id not found in context", "user_id not found in request context", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/time-slots/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logAndRespondError(w, http.StatusBadRequest, "invalid time slot id", "failed to parse time slot id", err)
		return
	}

	err = h.timeSlotService.DeleteTimeSlot(id, userID)
	if err != nil {
		if errors.Is(err, service.ErrTimeSlotNotFound) {
			logAndRespondError(w, http.StatusNotFound, "time slot not found", "time slot not found during delete", err)
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			logAndRespondError(w, http.StatusForbidden, "you can only manage your own time slots", "forbidden time slot delete attempt", err)
			return
		}
		logAndRespondError(w, http.StatusInternalServerError, "failed to delete time slot", "failed to delete time slot", err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "time slot deleted"})
}