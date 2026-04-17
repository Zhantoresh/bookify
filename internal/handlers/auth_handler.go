package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/usecase"
)

type AuthHandler struct {
	usecase *usecase.UserUsecase
}

func NewAuthHandler(u *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{usecase: u}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Name     string      `json:"name"`
		Role     domain.Role `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("failed to decode register request",
			slog.String("error", err.Error()),
		)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.Register(input.Email, input.Password, input.Name, input.Role)
	if err != nil {
		slog.Error("failed to register user",
			slog.String("email", input.Email),
			slog.String("error", err.Error()),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("failed to decode login request",
			slog.String("error", err.Error()),
		)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.usecase.Login(input.Email, input.Password)
	if err != nil {
		slog.Error("failed login attempt",
			slog.String("email", input.Email),
			slog.String("error", err.Error()),
		)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}