package handlers

import (
	"encoding/json"
	"net/http"
)


func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	
	response := map[string]string{
		"status":  "success",
		"message": "Успешный вход! Добро пожаловать в панель администратора",
		"role":    "admin",
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	
	json.NewEncoder(w).Encode(response)
}
