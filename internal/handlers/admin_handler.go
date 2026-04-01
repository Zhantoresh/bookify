package handlers

import (
	"encoding/json"
	"net/http"
)

// функция для проверки доступа админа
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	//ответ в формате JSON
	response := map[string]string{
		"status":  "success",
		"message": "Успешный вход! Добро пожаловать в панель администратора 🛡️",
		"role":    "admin",
	}

	// отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// oтправляем ответ
	json.NewEncoder(w).Encode(response)
}
