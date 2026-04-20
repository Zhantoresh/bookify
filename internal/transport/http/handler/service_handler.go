package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	"github.com/bookify/internal/service"
	"github.com/bookify/internal/transport/http/middleware"
)

type authWrapper func(http.Handler) http.Handler
type roleWrapper func(http.Handler, ...string) http.Handler

type ServiceHandler struct {
	services service.ServiceService
}

func NewServiceHandler(services service.ServiceService) *ServiceHandler {
	return &ServiceHandler{services: services}
}

func (h *ServiceHandler) HandleCollection(withRole roleWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.List(w, r)
		case http.MethodPost:
			withRole(http.HandlerFunc(h.Create), string(domain.RoleProvider), string(domain.RoleAdmin)).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (h *ServiceHandler) HandleByID(protected authWrapper, withRole roleWrapper) http.HandlerFunc {
	_ = protected
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/services/")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.GetByID(w, r, id)
		case http.MethodPut:
			withRole(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { h.Update(w, r, id) }), string(domain.RoleProvider), string(domain.RoleAdmin)).ServeHTTP(w, r)
		case http.MethodPatch:
			withRole(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { h.Patch(w, r, id) }), string(domain.RoleProvider), string(domain.RoleAdmin)).ServeHTTP(w, r)
		case http.MethodDelete:
			withRole(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { h.Delete(w, r, id) }), string(domain.RoleProvider), string(domain.RoleAdmin)).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (h *ServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := repository.ServiceFilter{
		Page:       parseInt(r.URL.Query().Get("page"), 1),
		Limit:      parseInt(r.URL.Query().Get("limit"), 20),
		ProviderID: r.URL.Query().Get("provider_id"),
		Search:     r.URL.Query().Get("search"),
	}
	if minPrice := r.URL.Query().Get("min_price"); minPrice != "" {
		if value, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = &value
		}
	}
	if maxPrice := r.URL.Query().Get("max_price"); maxPrice != "" {
		if value, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = &value
		}
	}
	items, pagination, err := h.services.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items, "pagination": pagination})
}

func (h *ServiceHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	items, pagination, err := h.services.ListMine(r.Context(), userID, parseInt(r.URL.Query().Get("page"), 1), parseInt(r.URL.Query().Get("limit"), 20))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items, "pagination": pagination})
}

func (h *ServiceHandler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.services.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request service.CreateServiceInput
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, err)
		return
	}
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role := domain.Role(r.Context().Value(middleware.RoleKey).(string))
	item, err := h.services.Create(r.Context(), userID, role, request)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	var request service.UpdateServiceInput
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, err)
		return
	}
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role := domain.Role(r.Context().Value(middleware.RoleKey).(string))
	item, err := h.services.Update(r.Context(), userID, role, id, request)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ServiceHandler) Patch(w http.ResponseWriter, r *http.Request, id string) {
	var request service.PatchServiceInput
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, err)
		return
	}
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role := domain.Role(r.Context().Value(middleware.RoleKey).(string))
	item, err := h.services.Patch(r.Context(), userID, role, id, request)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role := domain.Role(r.Context().Value(middleware.RoleKey).(string))
	if err := h.services.Delete(r.Context(), userID, role, id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	if n, err := strconv.Atoi(value); err == nil {
		return n
	}
	return fallback
}
