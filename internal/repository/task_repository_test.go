package repository

import (
	"context"
	"testing"

	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "Create and Get",
			run: func(t *testing.T) {
				repo := NewTaskRepository()
				task := models.Task{ID: "1", Title: "test", Description: "desc"}
				err := repo.Create(ctx, task)
				assert.NoError(t, err)
				got, ok, err := repo.Get(ctx, "1")
				assert.NoError(t, err)
				assert.True(t, ok)
				assert.Equal(t, task, got)
			},
		},
		{
			name: "Update",
			run: func(t *testing.T) {
				repo := NewTaskRepository()
				task := models.Task{ID: "1", Title: "test", Description: "desc"}
				_ = repo.Create(ctx, task)
				task.Title = "updated"
				err := repo.Update(ctx, task)
				assert.NoError(t, err)
				got, ok, err := repo.Get(ctx, "1")
				assert.NoError(t, err)
				assert.True(t, ok)
				assert.Equal(t, "updated", got.Title)
			},
		},
		{
			name: "GetAll",
			run: func(t *testing.T) {
				repo := NewTaskRepository()
				task := models.Task{ID: "1", Title: "test", Description: "desc"}
				_ = repo.Create(ctx, task)
				all := repo.GetAll(ctx)
				assert.Len(t, all, 1)
				assert.Equal(t, "test", all[0].Title)
			},
		},
		{
			name: "Delete",
			run: func(t *testing.T) {
				repo := NewTaskRepository()
				task := models.Task{ID: "1", Title: "test", Description: "desc"}
				_ = repo.Create(ctx, task)
				err := repo.Delete(ctx, "1")
				assert.NoError(t, err)
				_, ok, _ := repo.Get(ctx, "1")
				assert.False(t, ok)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
