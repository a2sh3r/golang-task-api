package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var body TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusBadRequest)
		logger.Log.Error("cant decode body", zap.Error(err))
		return
	}

	id, err := h.taskService.CreateTask(r.Context(), body.Title, body.Description)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		logger.Log.Error("error while creating task", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(id)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "task_id")
	task, exists, err := h.taskService.GetTask(r.Context(), id)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		logger.Log.Error("error while getting task", zap.Error(err))
		return
	}

	if !exists {
		http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
		logger.Log.Error("task does not exist", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "task_id")
	err := h.taskService.DeleteTask(r.Context(), id)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		logger.Log.Error("error while deleting task", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted"})
}
