package server

import (
	"net/http"

	"github.com/a2sh3r/golang-task-api.git/internal/middleware"
	"github.com/a2sh3r/golang-task-api.git/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	taskService service.TaskService
}

func NewHandler(taskService service.TaskService) *Handler {
	return &Handler{
		taskService: taskService,
	}
}

func NewRouter(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NewLoggingMiddleware())
	r.Use(middleware.NewGzipMiddleware())

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowerd", http.StatusMethodNotAllowed)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
	})

	r.Route("/api/tasks", func(r chi.Router) {
		r.Post("/", handler.CreateTask)
		r.Get("/{task_id}", handler.GetTask)
		r.Delete("/{task_id}", handler.DeleteTask)
	})

	return r
}
