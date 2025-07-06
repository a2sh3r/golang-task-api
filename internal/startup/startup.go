package startup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/config"
	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"github.com/a2sh3r/golang-task-api.git/internal/repository"
	"github.com/a2sh3r/golang-task-api.git/internal/server"
	"github.com/a2sh3r/golang-task-api.git/internal/service"
	"go.uber.org/zap"
)

type App struct {
	server   *http.Server
	taskRepo repository.TaskRepository
}

func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := logger.Initialize("debug"); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	logger.Log.Info("Server is starting",
		zap.String("address", cfg.RunAddress))

	taskRepo := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepo)

	handler := server.NewHandler(taskService)

	r := server.NewRouter(handler)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	return &App{
		server:   server,
		taskRepo: taskRepo,
	}, nil
}

func (a *App) Run(parentCtx context.Context) error {
	serverErrCh := make(chan error, 1)

	go func() {
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	select {
	case <-parentCtx.Done():
		return nil
	case err := <-serverErrCh:
		return err
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logger.Log.Info("shutting down server...")
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("server shutdown failed", zap.Error(err))
		return err
	}

	return nil
}
