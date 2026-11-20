package main

import (
	"log"
	"net/http"
	"time"

	"septikon-online/internal/httpapi"
	"septikon-online/internal/room"
)

func main() {
	// Создаём менеджер игровых комнат.
	roomManager := room.NewManager()

	// Создаём HTTP API.
	apiHandler := httpapi.NewHandler(roomManager)

	// Настраиваем маршруты.
	mux := http.NewServeMux()

	mux.Handle("/api/rooms", apiHandler)
	mux.Handle("/api/rooms/", apiHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Запускаем сервер.
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Septikon Online server started on :8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
