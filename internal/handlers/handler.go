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
	mux.HandleFunc("/specialists", h.handleSpecialists)
	mux.HandleFunc("/specialistsWithSlots/", h.handleSpecialistByID)
	mux.HandleFunc("/bookings", h.handleBookings)
}

func (h *Handler) handleSpecialists(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) handleSpecialistByID(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) handleBookings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createBooking(w, r)
	case http.MethodGet:
		h.getBookings(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
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

	const userID = 1 // Hardcoded user ID for MVP

	booking, err := h.bookingService.CreateBooking(userID, req.TimeSlotID)
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

	respondJSON(w, http.StatusCreated, booking)
}

func (h *Handler) getBookings(w http.ResponseWriter, r *http.Request) {
	const userID = 1 // Hardcoded user ID for MVP

	bookings, err := h.bookingService.GetUserBookings(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch bookings")
		return
	}

	respondJSON(w, http.StatusOK, bookings)
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
