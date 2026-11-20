package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Ответ сервера
type ServerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Проверка работоспособности сервера
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := ServerResponse{
		Status:  "ok",
		Message: "Septikon Online Backend работает!",
	}

	json.NewEncoder(w).Encode(response)
}

// Запуск сервера
func main() {
	http.HandleFunc("/health", healthHandler)

	log.Println("Septikon Online Backend запускается...")
	log.Println("Сервер доступен: http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
