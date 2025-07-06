package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/a2sh3r/golang-task-api.git/internal/repository"
	"go.uber.org/zap"
)

type TaskService interface {
	CreateTask(ctx context.Context, title string, description string) (string, error)
	GetTask(ctx context.Context, id string) (models.Task, bool, error)
	DeleteTask(ctx context.Context, id string) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{
		repo: repo,
	}
}

func (s *taskService) CreateTask(ctx context.Context, title string, description string) (string, error) {
	id := generateUniqueId()
	newTask := models.Task{
		ID:          id,
		Status:      models.Pending,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, newTask); err != nil {
		logger.Log.Error("failed to create task", zap.Error(err))
		return "", err
	}

	// Имитация обработки задачи
	go func() {
		duration := time.Duration(3+rand.Intn(3)) * time.Second
		time.Sleep(duration)

		if rand.Float64() < 0.2 {
			newTask.Status = models.Failed
		} else {
			newTask.Status = models.Completed
		}

		newTask.Duration = duration
		if err := s.repo.Update(context.Background(), newTask); err != nil {
			logger.Log.Error("error while updating task status", zap.Error(err))
		}
	}()

	return id, nil
}

func (s *taskService) GetTask(ctx context.Context, id string) (models.Task, bool, error) {
	return s.repo.Get(ctx, id)
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func generateUniqueId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
