package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/bookify/internal/service"
)

type Handler struct {
	specialistService *service.SpecialistService
	bookingService    *service.BookingService
}

func NewHandler(
	specialistService *service.SpecialistService,
	bookingService *service.BookingService,
) *Handler {
	return &Handler{
		specialistService: specialistService,
		bookingService:    bookingService,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/specialists", h.GetSpecialists)
	mux.HandleFunc("/specialistsWithSlots/", h.GetSpecialistByID)
	mux.HandleFunc("/bookings/", h.HandleBookingsByID)
	mux.HandleFunc("/bookings", h.HandleBookings)
}

func (h *Handler) GetSpecialists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logAndRespondError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed for /specialists", nil)
		return
	}

	specialists, err := h.specialistService.GetAllSpecialists()
	if err != nil {
		logAndRespondError(w, http.StatusInternalServerError, "failed to fetch specialists", "failed to fetch specialists", err)
		return
	}

	respondJSON(w, http.StatusOK, specialists)
}

func (h *Handler) GetSpecialistByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logAndRespondError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed for /specialistsWithSlots/", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/specialistsWithSlots/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logAndRespondError(w, http.StatusBadRequest, "invalid specialist id", "failed to parse specialist id", err)
		return
	}

	specialist, err := h.specialistService.GetSpecialistWithSlots(id)
	if err != nil {
		logAndRespondError(w, http.StatusNotFound, "specialist not found", "failed to fetch specialist with slots", err)
		return
	}

	respondJSON(w, http.StatusOK, specialist)
}

func (h *Handler) HandleBookings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createBooking(w, r)
	case http.MethodGet:
		h.getBookings(w, r)
	default:
		logAndRespondError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed for /bookings", nil)
	}
}

func (h *Handler) HandleBookingsByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		logAndRespondError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed for /bookings/{id}", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/bookings/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logAndRespondError(w, http.StatusBadRequest, "invalid booking id", "failed to parse booking id", err)
		return
	}

	h.cancelBooking(w, r, id)
}

func (h *Handler) createBooking(w http.ResponseWriter, r *http.Request) {
	type CreateBookingRequest struct {
		TimeSlotID int `json:"time_slot_id"`
	}

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logAndRespondError(w, http.StatusBadRequest, "invalid request body", "failed to decode create booking request", err)
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		logAndRespondError(w, http.StatusUnauthorized, "unauthorized", "user_id not found in request context", nil)
		return
	}

	bookingResponse, err := h.bookingService.CreateBookingWithDetails(userID, req.TimeSlotID)
	if err != nil {
		if errors.Is(err, service.ErrBookingAlreadyBooked) {
			logAndRespondError(w, http.StatusConflict, "this slot is already booked", "booking conflict: slot already booked", err)
			return
		}
		if errors.Is(err, service.ErrTimeSlotNotFound) {
			logAndRespondError(w, http.StatusNotFound, "time slot not found", "time slot not found during booking creation", err)
			return
		}
		logAndRespondError(w, http.StatusInternalServerError, "failed to create booking", "failed to create booking", err)
		return
	}

	respondJSON(w, http.StatusCreated, bookingResponse)
}

func (h *Handler) getBookings(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		logAndRespondError(w, http.StatusUnauthorized, "unauthorized", "user_id not found in request context", nil)
		return
	}

	bookings, err := h.bookingService.GetUserBookings(userID)
	if err != nil {
		logAndRespondError(w, http.StatusInternalServerError, "failed to fetch bookings", "failed to fetch user bookings", err)
		return
	}

	respondJSON(w, http.StatusOK, bookings)
}

func (h *Handler) cancelBooking(w http.ResponseWriter, r *http.Request, bookingID int) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		logAndRespondError(w, http.StatusUnauthorized, "unauthorized", "user_id not found in request context", nil)
		return
	}

	err := h.bookingService.CancelBooking(userID, bookingID)
	if err != nil {
		if errors.Is(err, service.ErrBookingNotFound) {
			logAndRespondError(w, http.StatusNotFound, "booking not found", "booking not found during cancel", err)
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			logAndRespondError(w, http.StatusForbidden, "you can only cancel your own bookings", "forbidden booking cancel attempt", err)
			return
		}
		if errors.Is(err, service.ErrBookingAlreadyCancelled) {
			logAndRespondError(w, http.StatusBadRequest, "booking is already cancelled", "booking already cancelled", err)
			return
		}
		logAndRespondError(w, http.StatusInternalServerError, "failed to cancel booking", "failed to cancel booking", err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "booking cancelled successfully"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func logAndRespondError(w http.ResponseWriter, status int, clientMessage string, logMessage string, err error) {
	attrs := []any{
		slog.Int("status", status),
		slog.String("client_message", clientMessage),
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	slog.Error(logMessage, attrs...)
	respondError(w, status, clientMessage)
}