package repository

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

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
	GetAll(ctx context.Context) []models.Task
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

func (r *taskRepository) GetAll(ctx context.Context) []models.Task {
	if err := checkStorageConsistency(r, ctx); err != nil {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	tasks := make([]models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks
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

func SaveTasksToFile(path string, tasks []models.Task) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logger.Log.Error("Error while opening storage file", zap.Error(err))
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Log.Error("Error closing storage file", zap.Error(err))
		}
	}()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(tasks); err != nil {
		logger.Log.Error("error while serializing storage data", zap.Error(err))
	}
}

func LoadTasksFromFile(path string) ([]models.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Log.Warn("Warning: file not found, starting with an empty task list", zap.String("file", path))
			return []models.Task{}, nil
		}
		logger.Log.Warn("Warning: could not open file, starting with an empty task list", zap.String("file", path), zap.Error(err))
		return []models.Task{}, nil
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Log.Error("Error closing storage file", zap.Error(err))
		}
	}()

	var tasks []models.Task
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		logger.Log.Warn("Warning: could not decode tasks from file, starting with an empty task list", zap.String("file", path), zap.Error(err))
		return []models.Task{}, nil
	}
	return tasks, nil
}

func StartAutoSave(ctx context.Context, path string, getTasks func() []models.Task, interval int) {
	if interval <= 0 {
		<-ctx.Done()
		SaveTasksToFile(path, getTasks())
		return
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			SaveTasksToFile(path, getTasks())
		case <-ctx.Done():
			SaveTasksToFile(path, getTasks())
			return
		}
	}
}
