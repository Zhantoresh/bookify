package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	"github.com/bookify/internal/service"
	"github.com/bookify/internal/transport/http/middleware"
	"github.com/bookify/internal/worker"
)

type AppointmentHandler struct {
	appointments service.AppointmentService
	workerPool   *worker.WorkerPool
	logger       *slog.Logger
}

func NewAppointmentHandler(appointments service.AppointmentService, workerPool *worker.WorkerPool, logger *slog.Logger) *AppointmentHandler {
	return &AppointmentHandler{appointments: appointments, workerPool: workerPool, logger: logger}
}

func (h *AppointmentHandler) HandleCollection(protected authWrapper, withRole roleWrapper) http.Handler {
	_ = withRole
	return protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			if r.Context().Value(middleware.RoleKey) != string(domain.RoleClient) {
				writeError(w, domain.ErrForbidden)
				return
			}
			h.Create(w, r)
		case http.MethodGet:
			if r.Context().Value(middleware.RoleKey) != string(domain.RoleAdmin) {
				writeError(w, domain.ErrForbidden)
				return
			}
			h.List(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func (h *AppointmentHandler) HandleByID(protected authWrapper) http.Handler {
	return protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/appointments/")
		if path == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		switch {
		case strings.HasSuffix(path, "/confirm") && r.Method == http.MethodPatch:
			id := strings.TrimSuffix(path, "/confirm")
			h.Confirm(w, r, strings.TrimSuffix(id, "/"))
		case strings.HasSuffix(path, "/cancel") && r.Method == http.MethodPatch:
			id := strings.TrimSuffix(path, "/cancel")
			h.Cancel(w, r, strings.TrimSuffix(id, "/"))
		case strings.HasSuffix(path, "/complete") && r.Method == http.MethodPatch:
			id := strings.TrimSuffix(path, "/complete")
			h.Complete(w, r, strings.TrimSuffix(id, "/"))
		case r.Method == http.MethodGet:
			h.GetByID(w, r, path)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ServiceID string `json:"service_id"`
		StartTime string `json:"start_time"`
		Notes     string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, err)
		return
	}
	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		writeError(w, err)
		return
	}
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	item, err := h.appointments.Create(r.Context(), userID, service.CreateAppointmentInput{
		ServiceID: request.ServiceID,
		StartTime: startTime,
		Notes:     request.Notes,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	h.workerPool.Submit(func() error {
		h.logger.Info("appointment_created", "appointment_id", item.ID)
		return nil
	})
	writeJSON(w, http.StatusCreated, item)
}

func (h *AppointmentHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := repository.AppointmentFilter{
		Page:       parseInt(r.URL.Query().Get("page"), 1),
		Limit:      parseInt(r.URL.Query().Get("limit"), 20),
		Status:     r.URL.Query().Get("status"),
		ClientID:   r.URL.Query().Get("client_id"),
		ProviderID: r.URL.Query().Get("provider_id"),
	}
	if fromDate := r.URL.Query().Get("from_date"); fromDate != "" {
		if value, err := parseDateFilter(fromDate, false); err == nil {
			filter.FromDate = &value
		}
	}
	if toDate := r.URL.Query().Get("to_date"); toDate != "" {
		if value, err := parseDateFilter(toDate, true); err == nil {
			filter.ToDate = &value
		}
	}
	items, pagination, err := h.appointments.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items, "pagination": pagination})
}

func (h *AppointmentHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role := domain.Role(r.Context().Value(middleware.RoleKey).(string))
	items, pagination, err := h.appointments.ListMine(r.Context(), userID, role, parseInt(r.URL.Query().Get("page"), 1), parseInt(r.URL.Query().Get("limit"), 20))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items, "pagination": pagination})
}

func (h *AppointmentHandler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role := domain.Role(r.Context().Value(middleware.RoleKey).(string))
	item, err := h.appointments.GetByID(r.Context(), userID, role, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *AppointmentHandler) Confirm(w http.ResponseWriter, r *http.Request, id string) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role := domain.Role(r.Context().Value(middleware.RoleKey).(string))
	item, err := h.appointments.Confirm(r.Context(), userID, role, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":           item.ID,
		"status":       item.Status,
		"confirmed_at": item.ConfirmedAt,
	})
}

func (h *AppointmentHandler) Cancel(w http.ResponseWriter, r *http.Request, id string) {
	var request struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&request)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role := domain.Role(r.Context().Value(middleware.RoleKey).(string))
	item, err := h.appointments.Cancel(r.Context(), userID, role, id, request.Reason)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":                  item.ID,
		"status":              item.Status,
		"cancelled_at":        item.CancelledAt,
		"cancellation_reason": item.CancellationReason,
	})
}

func (h *AppointmentHandler) Complete(w http.ResponseWriter, r *http.Request, id string) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role := domain.Role(r.Context().Value(middleware.RoleKey).(string))
	item, err := h.appointments.Complete(r.Context(), userID, role, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *AppointmentHandler) AvailableSlots(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Query().Get("service_id")
	dateValue := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateValue)
	if err != nil {
		writeError(w, err)
		return
	}
	slots, err := h.appointments.AvailableSlots(r.Context(), serviceID, date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"service_id":      serviceID,
		"date":            date.Format("2006-01-02"),
		"available_slots": slots,
	})
}

func parseDateFilter(value string, endOfDay bool) (time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, err
	}
	if endOfDay {
		return parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second), nil
	}
	return parsed, nil
}
