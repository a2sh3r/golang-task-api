package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/a2sh3r/golang-task-api.git/internal/models"
)

var (
	ErrStorageNil        = errors.New("storage is nil")
	ErrTaskMapNil        = errors.New("task map is nil")
	ErrInvalidTask       = errors.New("invalid input task")
	ErrTaskAlreadyExists = errors.New("task with this ID already exists")
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

func NewtaskRepository() TaskRepository {
	return &taskRepository{
		tasks: make(map[string]models.Task),
	}
}

func (r *taskRepository) Create(ctx context.Context, task models.Task) error {
	if err := checkStorageConsistency(r, ctx); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return ErrTaskAlreadyExists
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *taskRepository) Get(ctx context.Context, id string) (models.Task, bool, error) {
	if err := checkStorageConsistency(r, ctx); err != nil {
		return models.Task{}, false, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	return task, exists, nil
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	if err := checkStorageConsistency(r, ctx); err != nil {
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
		return errors.New("task not found")
	}

	r.tasks[task.ID] = task
	return nil
}

func checkStorageConsistency(r *taskRepository, ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if r.tasks == nil {
		return ErrTaskMapNil
	}

	return nil
}
