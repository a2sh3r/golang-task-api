package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/config"
	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"github.com/a2sh3r/golang-task-api.git/internal/repository"
	"github.com/a2sh3r/golang-task-api.git/internal/server"
	"github.com/a2sh3r/golang-task-api.git/internal/service"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	if err := logger.Initialize("debug"); err != nil {
		panic(err)
	}

	taskRepo := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepo)

	handler := server.NewHandler(taskService)

	r := server.NewRouter(handler)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	<-ctx.Done()
	err = <-serverErrCh
	if err != nil {
		logger.Log.Fatal("server error", zap.Error(err))
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	logger.Log.Info("shutting down server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("server shutdown failed", zap.Error(err))
		panic(err)
	}

	logger.Log.Info("server exited gracefully")
}
