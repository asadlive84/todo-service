package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"todo-service/internal/usecase"

	"github.com/gorilla/mux"
)

type Handler struct {
	todoUC *usecase.TodoUseCase
	fileUC *usecase.FileUseCase
}

func NewHandler(todoUC *usecase.TodoUseCase, fileUC *usecase.FileUseCase) *Handler {
	return &Handler{todoUC: todoUC, fileUC: fileUC}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"time":   time.Now(),
	})
}

func RegisterRoutes(r *mux.Router, h *Handler) {

	r.HandleFunc("/todo", h.CreateTodo).Methods("POST")

	// File route
	r.HandleFunc("/upload", h.Upload).Methods("POST")

	// Health route
	r.HandleFunc("/health", h.Health).Methods("GET")
}
