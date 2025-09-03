package startup

import (
	"context"
	"fmt"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/config"
	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"github.com/a2sh3r/golang-task-api.git/internal/repository/postgres"
	"github.com/a2sh3r/golang-task-api.git/internal/server"
	"github.com/a2sh3r/golang-task-api.git/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	app         *fiber.App
	taskRepo    *postgres.TaskRepository
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

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	taskRepo := postgres.NewTaskRepository(dbpool)
	taskService := service.NewTaskService(taskRepo)
	handler := server.NewHandler(taskService)
	app := server.NewRouter(handler)

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		app:         app,
		taskRepo:    taskRepo,
		shutdownCtx: ctx,
		cancelFunc:  cancel,
		cfg:         cfg,
	}, nil
}

func (a *App) Run(parentCtx context.Context) error {
	serverErrCh := make(chan error, 1)

	go func() {
		err := a.app.Listen(a.cfg.RunAddress)
		if err != nil {
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

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := a.app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Log.Error("server shutdown failed", zap.Error(err))
		return err
	}

	return nil
}
