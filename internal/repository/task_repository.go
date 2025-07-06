package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"go.uber.org/zap"
)

var (
	ErrStorageNil        = errors.New("storage is nil")
	ErrTaskMapNil        = errors.New("task map is nil")
	ErrInvalidTask       = errors.New("invalid input task")
	ErrTaskAlreadyExists = errors.New("task with this ID already exists")
	ErrTaskNotFound      = errors.New("task not found")
)

type TaskRepository interface {
	Create(ctx context.Context, task models.Task) error
	Get(ctx context.Context, id string) (models.Task, bool, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, task models.Task) error
}

type taskRepository struct {
	mu    sync.Mutex
	tasks map[string]models.Task
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{
		tasks: make(map[string]models.Task),
	}
}

func (r *taskRepository) Create(ctx context.Context, task models.Task) error {
	if err := checkStorageConsistency(r, ctx); err != nil {
		logger.Log.Error("storage is not consistent", zap.Error(err))
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		logger.Log.Error("task already exists", zap.Error(ErrTaskAlreadyExists))
		return ErrTaskAlreadyExists
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *taskRepository) Get(ctx context.Context, id string) (models.Task, bool, error) {
	if err := checkStorageConsistency(r, ctx); err != nil {
		logger.Log.Error("storage is not consistent", zap.Error(err))
		return models.Task{}, false, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	return task, exists, nil
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	if err := checkStorageConsistency(r, ctx); err != nil {
		logger.Log.Error("storage is not consistent", zap.Error(err))
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tasks, id)

	return nil
}

func (r *taskRepository) Update(ctx context.Context, task models.Task) error {
	if err := checkStorageConsistency(r, ctx); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		logger.Log.Error("task not found", zap.Error(ErrTaskNotFound))
		return ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

func checkStorageConsistency(r *taskRepository, ctx context.Context) error {
	if ctx.Err() != nil {
		logger.Log.Error("storage context error", zap.Error(ctx.Err()))
		return ctx.Err()
	}

	if r.tasks == nil {
		logger.Log.Error("storage is nil", zap.Error(ErrTaskMapNil))
		return ErrTaskMapNil
	}

	return nil
}
