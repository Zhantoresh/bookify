package handler

import (
	"net/http"

	"github.com/bookify/internal/service"
)

type AdminHandler struct {
	admin service.AdminService
}

func NewAdminHandler(admin service.AdminService) *AdminHandler {
	return &AdminHandler{admin: admin}
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := h.admin.Dashboard(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}
