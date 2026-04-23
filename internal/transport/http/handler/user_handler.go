package handler

import (
	"net/http"
	"strings"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	"github.com/bookify/internal/service"
	"github.com/bookify/internal/transport/http/middleware"
)

type UserHandler struct {
	users service.UserService
	admin service.AdminService
}

func NewUserHandler(users service.UserService, admin service.AdminService) *UserHandler {
	return &UserHandler{users: users, admin: admin}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) AdminCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	filter := repository.UserFilter{Role: domain.Role(r.URL.Query().Get("role"))}
	if filter.Role != "" && filter.Role != domain.RoleAdmin && filter.Role != domain.RoleClient && filter.Role != domain.RoleProvider {
		writeError(w, domain.ErrValidation)
		return
	}
	users, err := h.admin.ListUsers(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": users})
}

func (h *UserHandler) AdminByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/users/")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		user, err := h.admin.GetUser(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, user)
	case http.MethodPatch:
		var request struct {
			Role string `json:"role"`
		}
		if err := decodeJSONBody(r, &request); err != nil {
			writeError(w, err)
			return
		}
		actorID, _ := r.Context().Value(middleware.UserIDKey).(string)
		user, err := h.admin.UpdateUserRole(r.Context(), actorID, id, domain.Role(request.Role))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, user)
	case http.MethodDelete:
		actorID, _ := r.Context().Value(middleware.UserIDKey).(string)
		if err := h.admin.DeleteUser(r.Context(), actorID, id); err != nil {
			writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
