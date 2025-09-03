package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, title, description string) (*models.Task, error) {
	args := m.Called(ctx, title, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskService) UpdateTask(ctx context.Context, id string, req models.UpdateTaskRequest) (*models.Task, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestApp() *fiber.App {
	app := fiber.New()
	return app
}

func TestHandler_CreateTask(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    models.CreateTaskRequest
		setupMock      func(*MockTaskService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful task creation",
			requestBody: models.CreateTaskRequest{
				Title:       "Test Task",
				Description: "Test Description",
			},
			setupMock: func(mockService *MockTaskService) {
				expectedTask := &models.Task{
					ID:          1,
					Title:       "Test Task",
					Description: "Test Description",
					Status:      models.New,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockService.On("CreateTask", mock.Anything, "Test Task", "Test Description").Return(expectedTask, nil)
			},
			expectedStatus: fiber.StatusCreated,
		},
		{
			name: "invalid request body",
			requestBody: models.CreateTaskRequest{
				Title:       "",
				Description: "Test Description",
			},
			setupMock: func(mockService *MockTaskService) {
				mockService.On("CreateTask", mock.Anything, "", "Test Description").Return(nil, assert.AnError)
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()
			mockService := new(MockTaskService)
			tt.setupMock(mockService)

			handler := &Handler{taskService: mockService}
			app.Post("/tasks", handler.CreateTask)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_GetAllTasks(t *testing.T) {
	mockService := new(MockTaskService)
	expectedTasks := []models.Task{
		{ID: 1, Title: "Task 1", Description: "Description 1", Status: models.New},
		{ID: 2, Title: "Task 2", Description: "Description 2", Status: models.InProgress},
	}
	mockService.On("GetAllTasks", mock.Anything).Return(expectedTasks, nil)

	app := setupTestApp()
	handler := &Handler{taskService: mockService}
	app.Get("/tasks", handler.GetAllTasks)

	req := httptest.NewRequest("GET", "/tasks", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestHandler_GetTask(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		setupMock      func(*MockTaskService)
		expectedStatus int
	}{
		{
			name:   "successful task retrieval",
			taskID: "1",
			setupMock: func(mockService *MockTaskService) {
				expectedTask := &models.Task{
					ID:          1,
					Title:       "Test Task",
					Description: "Test Description",
					Status:      models.New,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockService.On("GetTask", mock.Anything, "1").Return(expectedTask, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:   "task not found",
			taskID: "999",
			setupMock: func(mockService *MockTaskService) {
				mockService.On("GetTask", mock.Anything, "999").Return(nil, assert.AnError)
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()
			mockService := new(MockTaskService)
			tt.setupMock(mockService)

			handler := &Handler{taskService: mockService}
			app.Get("/tasks/:id", handler.GetTask)

			req := httptest.NewRequest("GET", "/tasks/"+tt.taskID, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_UpdateTask(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		requestBody    models.UpdateTaskRequest
		setupMock      func(*MockTaskService)
		expectedStatus int
	}{
		{
			name:   "successful task update",
			taskID: "1",
			requestBody: models.UpdateTaskRequest{
				Title:  stringPtr("Updated Title"),
				Status: statusPtr(models.InProgress),
			},
			setupMock: func(mockService *MockTaskService) {
				expectedTask := &models.Task{
					ID:          1,
					Title:       "Updated Title",
					Description: "Original Description",
					Status:      models.InProgress,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockService.On("UpdateTask", mock.Anything, "1", mock.AnythingOfType("models.UpdateTaskRequest")).Return(expectedTask, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:   "task not found",
			taskID: "999",
			requestBody: models.UpdateTaskRequest{
				Title: stringPtr("Updated Title"),
			},
			setupMock: func(mockService *MockTaskService) {
				mockService.On("UpdateTask", mock.Anything, "999", mock.AnythingOfType("models.UpdateTaskRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()
			mockService := new(MockTaskService)
			tt.setupMock(mockService)

			handler := &Handler{taskService: mockService}
			app.Put("/tasks/:id", handler.UpdateTask)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/tasks/"+tt.taskID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteTask(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		setupMock      func(*MockTaskService)
		expectedStatus int
	}{
		{
			name:   "successful task deletion",
			taskID: "1",
			setupMock: func(mockService *MockTaskService) {
				mockService.On("DeleteTask", mock.Anything, "1").Return(nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:   "task not found",
			taskID: "999",
			setupMock: func(mockService *MockTaskService) {
				mockService.On("DeleteTask", mock.Anything, "999").Return(assert.AnError)
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()
			mockService := new(MockTaskService)
			tt.setupMock(mockService)

			handler := &Handler{taskService: mockService}
			app.Delete("/tasks/:id", handler.DeleteTask)

			req := httptest.NewRequest("DELETE", "/tasks/"+tt.taskID, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func statusPtr(s models.Status) *models.Status {
	return &s
}
