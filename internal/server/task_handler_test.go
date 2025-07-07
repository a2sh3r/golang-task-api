package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/a2sh3r/golang-task-api.git/internal/mocks/service_mocks"
	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateTask(t *testing.T) {
	tests := []struct {
		name      string
		body      map[string]string
		mockSetup func(s *mocks.MockTaskService)
		wantCode  int
	}{
		{
			name: "success",
			body: map[string]string{"title": "t1", "description": "d1"},
			mockSetup: func(s *mocks.MockTaskService) {
				s.EXPECT().
					CreateTask(gomock.Any(), "t1", "d1").
					Return("id1", nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "service error",
			body: map[string]string{"title": "t2", "description": "d2"},
			mockSetup: func(s *mocks.MockTaskService) {
				s.EXPECT().
					CreateTask(gomock.Any(), "t2", "d2").
					Return("", assert.AnError)
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := mocks.NewMockTaskService(ctrl)
			tt.mockSetup(svc)
			h := NewHandler(svc)
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/tasks/", bytes.NewReader(body))
			w := httptest.NewRecorder()
			h.CreateTask(w, req)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestHandler_GetTask(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockSetup func(s *mocks.MockTaskService)
		wantCode  int
	}{
		{
			name: "found",
			id:   "id1",
			mockSetup: func(s *mocks.MockTaskService) {
				s.EXPECT().
					GetTask(gomock.Any(), "id1").
					Return(models.Task{ID: "id1"}, true, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "not found",
			id:   "id2",
			mockSetup: func(s *mocks.MockTaskService) {
				s.EXPECT().
					GetTask(gomock.Any(), "id2").
					Return(models.Task{}, false, nil)
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "service error",
			id:   "id3",
			mockSetup: func(s *mocks.MockTaskService) {
				s.EXPECT().
					GetTask(gomock.Any(), "id3").
					Return(models.Task{}, false, errors.New("fail"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := mocks.NewMockTaskService(ctrl)
			tt.mockSetup(svc)
			h := NewHandler(svc)

			r := chi.NewRouter()
			r.Get("/api/tasks/{task_id}", h.GetTask)

			req := httptest.NewRequest(http.MethodGet, "/api/tasks/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
