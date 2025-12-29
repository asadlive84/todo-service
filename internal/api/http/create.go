package http

import (
	"encoding/json"
	"net/http"
	"time"
	"todo-service/internal/domain/entity"

	"github.com/rs/zerolog/log"
)

func (h *Http) CreateTodo(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("CreateTodo handler called")
	var req struct {
		Description string `json:"description"`
		DueDate     string `json:"dueDate"`
		FileID      string `json:"fileId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.DueDate == "" {
		req.DueDate = time.Now().String()
	}
	due, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		http.Error(w, "invalid dueDate", http.StatusBadRequest)
		return
	}

	if req.FileID != "" {
		isExists, err := h.fileUC.ValidateFileID(r.Context(), req.FileID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !isExists {
			http.Error(w, "fileId does not exist", http.StatusBadRequest)
			return
		}
	}

	en := &entity.TodoItem{
		Description: req.Description,
		DueDate:     due,
		FileID:      req.FileID,
	}

	err = h.todoUC.Create(r.Context(), en)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(en)
}
