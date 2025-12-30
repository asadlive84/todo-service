package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"todo-service/internal/port"
)

type Http struct {
	todoUC port.TodoUseCasePort
	fileUC port.FileUseCasePort
}

func NewHandler(todoUC port.TodoUseCasePort, fileUC port.FileUseCasePort) *Http {
	return &Http{todoUC: todoUC, fileUC: fileUC}
}

func (h *Http) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"time":   time.Now(),
	})
}

func RegisterRoutes(r *mux.Router, h *Http) {

	r.HandleFunc("/todo", h.CreateTodo).Methods("POST")

	// File route
	r.HandleFunc("/upload", h.Upload).Methods("POST")

	// Health route
	r.HandleFunc("/health", h.Health).Methods("GET")
}
