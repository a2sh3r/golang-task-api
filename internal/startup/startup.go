package startup

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/a2sh3r/golang-task-api.git/internal/config"
	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/a2sh3r/golang-task-api.git/internal/repository"
	"github.com/a2sh3r/golang-task-api.git/internal/server"
	"github.com/a2sh3r/golang-task-api.git/internal/service"
	"go.uber.org/zap"
)

type App struct {
	server      *http.Server
	taskRepo    repository.TaskRepository
	shutdownCtx context.Context
	cancelFunc  context.CancelFunc
	cfg         *config.Config
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

	tasks, err := repository.LoadTasksFromFile(cfg.FileStoragePath)
	if err != nil {
		logger.Log.Warn("Could not load tasks from file, starting fresh", zap.Error(err))
	}

	taskRepo := repository.NewTaskRepository()
	for _, task := range tasks {
		if err := taskRepo.Create(context.Background(), task); err != nil {
			logger.Log.Error("Failed to restore task", zap.String("task_id", task.ID), zap.Error(err))
		}
	}

	taskService := service.NewTaskService(taskRepo)

	handler := server.NewHandler(taskService)

	r := server.NewRouter(handler)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		server:      server,
		taskRepo:    taskRepo,
		shutdownCtx: ctx,
		cancelFunc:  cancel,
		cfg:         cfg,
	}

	go repository.StartAutoSave(
		app.shutdownCtx,
		cfg.FileStoragePath,
		func() []models.Task {
			return app.taskRepo.GetAll(context.Background())
		},
		cfg.StoreInterval,
	)

	return app, nil
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
	a.cancelFunc()

	logger.Log.Info("shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Log.Error("server shutdown failed", zap.Error(err))
		return err
	}

	repository.SaveTasksToFile(
		a.cfg.FileStoragePath,
		a.taskRepo.GetAll(context.Background()),
	)

	if err := a.server.Shutdown(ctx); err != nil {
		logger.Log.Error("server shutdown failed", zap.Error(err))
		return err
	}

	return nil
}
