package handlers

import (
	"bookify/internal/domain"
	"bookify/internal/usecase"
	"encoding/json"
	"net/http"
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
		Role     domain.Role `json:"role"` 
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.Register(input.Email, input.Password, input.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.usecase.Login(input.Email, input.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}