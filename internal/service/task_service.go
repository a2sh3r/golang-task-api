package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task models.Task) error
	GetByID(ctx context.Context, id int) (*models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	Update(ctx context.Context, id int, task models.Task) error
	Delete(ctx context.Context, id int) error
}

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, title, description string) (*models.Task, error)
	GetTask(ctx context.Context, idStr string) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	UpdateTask(ctx context.Context, idStr string, req models.UpdateTaskRequest) (*models.Task, error)
	DeleteTask(ctx context.Context, idStr string) error
}

type TaskService struct {
	taskRepo TaskRepository
}

func NewTaskService(taskRepo TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description string) (*models.Task, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	now := time.Now()
	task := models.Task{
		Title:       title,
		Description: description,
		Status:      models.New,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &task, nil
}

func (s *TaskService) GetTask(ctx context.Context, idStr string) (*models.Task, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %s", idStr)
	}

	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.taskRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, idStr string, req models.UpdateTaskRequest) (*models.Task, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %s", idStr)
	}

	existingTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if existingTask == nil {
		return nil, fmt.Errorf("task not found")
	}

	if req.Title != nil {
		existingTask.Title = *req.Title
	}
	if req.Description != nil {
		existingTask.Description = *req.Description
	}
	if req.Status != nil {
		existingTask.Status = *req.Status
	}

	existingTask.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, id, *existingTask); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return existingTask, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid task id: %s", idStr)
	}

	if err := s.taskRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
