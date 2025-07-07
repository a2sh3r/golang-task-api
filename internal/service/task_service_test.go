package service

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/a2sh3r/golang-task-api.git/internal/mocks/repository_mocks"
	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskService_CreateTask(t *testing.T) {
	type args struct {
		title, description string
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(repo *mocks.MockTaskRepository)
		wantErr   bool
	}{
		{
			name: "success",
			args: args{"title1", "desc1"},
			mockSetup: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(models.Task{})).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			args: args{"title2", "desc2"},
			mockSetup: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(models.Task{})).
					Return(errors.New("fail"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mocks.NewMockTaskRepository(ctrl)
			tt.mockSetup(repo)
			svc := NewTaskService(repo)
			_, err := svc.CreateTask(context.Background(), tt.args.title, tt.args.description)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskService_GetTask(t *testing.T) {
	type args struct{ id string }
	tests := []struct {
		name      string
		args      args
		mockSetup func(repo *mocks.MockTaskRepository)
		wantOk    bool
		wantErr   bool
	}{
		{
			name: "found",
			args: args{"id1"},
			mockSetup: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					Get(gomock.Any(), "id1").
					Return(models.Task{ID: "id1"}, true, nil)
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "not found",
			args: args{"id2"},
			mockSetup: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					Get(gomock.Any(), "id2").
					Return(models.Task{}, false, nil)
			},
			wantOk:  false,
			wantErr: false,
		},
		{
			name: "repo error",
			args: args{"id3"},
			mockSetup: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					Get(gomock.Any(), "id3").
					Return(models.Task{}, false, errors.New("fail"))
			},
			wantOk:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mocks.NewMockTaskRepository(ctrl)
			tt.mockSetup(repo)
			svc := NewTaskService(repo)
			_, ok, err := svc.GetTask(context.Background(), tt.args.id)
			assert.Equal(t, tt.wantOk, ok)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	type args struct{ id string }
	tests := []struct {
		name      string
		args      args
		mockSetup func(repo *mocks.MockTaskRepository)
		wantErr   bool
	}{
		{
			name: "success",
			args: args{"id1"},
			mockSetup: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					Delete(gomock.Any(), "id1").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			args: args{"id2"},
			mockSetup: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					Delete(gomock.Any(), "id2").
					Return(errors.New("fail"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mocks.NewMockTaskRepository(ctrl)
			tt.mockSetup(repo)
			svc := NewTaskService(repo)
			err := svc.DeleteTask(context.Background(), tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
