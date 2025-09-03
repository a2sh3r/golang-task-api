package service

import (
	"context"
	"testing"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, id int, task models.Task) error {
	args := m.Called(ctx, id, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTaskService_CreateTask(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		setupMock   func(*MockTaskRepository)
		wantErr     bool
	}{
		{
			name:        "successful task creation",
			title:       "Test Task",
			description: "Test Description",
			setupMock: func(mockRepo *MockTaskRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("models.Task")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "empty title should fail",
			title:       "",
			description: "Test Description",
			setupMock:   func(mockRepo *MockTaskRepository) {},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.setupMock(mockRepo)

			service := NewTaskService(mockRepo)
			task, err := service.CreateTask(context.Background(), tt.title, tt.description)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, tt.title, task.Title)
				assert.Equal(t, tt.description, task.Description)
				assert.Equal(t, models.New, task.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_GetTask(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*MockTaskRepository)
		wantErr   bool
		wantTask  *models.Task
	}{
		{
			name: "successful task retrieval",
			id:   "1",
			setupMock: func(mockRepo *MockTaskRepository) {
				expectedTask := &models.Task{
					ID:          1,
					Title:       "Test Task",
					Description: "Test Description",
					Status:      models.New,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.On("GetByID", mock.Anything, 1).Return(expectedTask, nil)
			},
			wantErr:  false,
			wantTask: &models.Task{ID: 1, Title: "Test Task", Description: "Test Description", Status: models.New},
		},
		{
			name:      "invalid id should fail",
			id:        "invalid",
			setupMock: func(mockRepo *MockTaskRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.setupMock(mockRepo)

			service := NewTaskService(mockRepo)
			task, err := service.GetTask(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, tt.wantTask.ID, task.ID)
				assert.Equal(t, tt.wantTask.Title, task.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_GetAllTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	expectedTasks := []models.Task{
		{ID: 1, Title: "Task 1", Description: "Description 1", Status: models.New},
		{ID: 2, Title: "Task 2", Description: "Description 2", Status: models.InProgress},
	}
	mockRepo.On("GetAll", mock.Anything).Return(expectedTasks, nil)

	service := NewTaskService(mockRepo)
	tasks, err := service.GetAllTasks(context.Background())

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, expectedTasks[0].Title, tasks[0].Title)
	assert.Equal(t, expectedTasks[1].Title, tasks[1].Title)

	mockRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTask(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		request   models.UpdateTaskRequest
		setupMock func(*MockTaskRepository)
		wantErr   bool
	}{
		{
			name: "successful task update",
			id:   "1",
			request: models.UpdateTaskRequest{
				Title:  stringPtr("Updated Title"),
				Status: statusPtr(models.InProgress),
			},
			setupMock: func(mockRepo *MockTaskRepository) {
				existingTask := &models.Task{
					ID:          1,
					Title:       "Original Title",
					Description: "Original Description",
					Status:      models.New,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.On("GetByID", mock.Anything, 1).Return(existingTask, nil)
				mockRepo.On("Update", mock.Anything, 1, mock.AnythingOfType("models.Task")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid id should fail",
			id:   "invalid",
			request: models.UpdateTaskRequest{
				Title: stringPtr("Updated Title"),
			},
			setupMock: func(mockRepo *MockTaskRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.setupMock(mockRepo)

			service := NewTaskService(mockRepo)
			task, err := service.UpdateTask(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*MockTaskRepository)
		wantErr   bool
	}{
		{
			name: "successful task deletion",
			id:   "1",
			setupMock: func(mockRepo *MockTaskRepository) {
				mockRepo.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "invalid id should fail",
			id:        "invalid",
			setupMock: func(mockRepo *MockTaskRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.setupMock(mockRepo)

			service := NewTaskService(mockRepo)
			err := service.DeleteTask(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func statusPtr(s models.Status) *models.Status {
	return &s
}
