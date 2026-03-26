package handlers

import (
	"encoding/json"
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
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	specialists, err := h.specialistService.GetAllSpecialists()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch specialists")
		return
	}

	respondJSON(w, http.StatusOK, specialists)
}

func (h *Handler) GetSpecialistByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/specialistsWithSlots/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid specialist id")
		return
	}

	specialist, err := h.specialistService.GetSpecialistWithSlots(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "specialist not found")
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
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) HandleBookingsByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/bookings/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid booking id")
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
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookingResponse, err := h.bookingService.CreateBookingWithDetails(userID, req.TimeSlotID)
	if err != nil {
		if err.Error() == "this slot is already booked" {
			respondError(w, http.StatusConflict, "this slot is already booked")
			return
		}
		if err.Error() == "time slot not found" {
			respondError(w, http.StatusNotFound, "time slot not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create booking")
		return
	}

	respondJSON(w, http.StatusCreated, bookingResponse)
}

func (h *Handler) getBookings(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookings, err := h.bookingService.GetUserBookings(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch bookings")
		return
	}

	respondJSON(w, http.StatusOK, bookings)
}

func (h *Handler) cancelBooking(w http.ResponseWriter, r *http.Request, bookingID int) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.bookingService.CancelBooking(userID, bookingID)
	if err != nil {
		if err.Error() == "booking not found" {
			respondError(w, http.StatusNotFound, "booking not found")
			return
		}
		if err.Error() == "forbidden" {
			respondError(w, http.StatusForbidden, "you can only cancel your own bookings")
			return
		}
		if err.Error() == "booking is already cancelled" {
			respondError(w, http.StatusBadRequest, "booking is already cancelled")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to cancel booking")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "booking cancelled successfully"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
