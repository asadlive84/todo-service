package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	port "todo-service/internal/domain/interface"
)

type Http struct {
	todoUC port.TodoUseCase
	fileUC port.FileUseCase
}

func NewHandler(todoUC port.TodoUseCase, fileUC port.FileUseCase) *Http {
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
