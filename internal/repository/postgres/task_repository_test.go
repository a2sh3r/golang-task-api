package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/golangtaskapi?sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		DROP TABLE IF EXISTS tasks;
		CREATE TABLE tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now()
		);
	`)
	require.NoError(t, err)

	return pool
}

func TestTaskRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	ctx := context.Background()

	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.New,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, task)
	assert.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM tasks").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestTaskRepository_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	ctx := context.Background()

	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.New,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, task)
	require.NoError(t, err)

	retrievedTask, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedTask)
	assert.Equal(t, "Test Task", retrievedTask.Title)
	assert.Equal(t, "Test Description", retrievedTask.Description)
	assert.Equal(t, models.New, retrievedTask.Status)
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	ctx := context.Background()

	task, err := repo.GetByID(ctx, 999)
	assert.NoError(t, err)
	assert.Nil(t, task)
}

func TestTaskRepository_GetAll(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	ctx := context.Background()

	task1 := models.Task{
		Title:       "Task 1",
		Description: "Description 1",
		Status:      models.New,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	task2 := models.Task{
		Title:       "Task 2",
		Description: "Description 2",
		Status:      models.InProgress,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, task1)
	require.NoError(t, err)
	err = repo.Create(ctx, task2)
	require.NoError(t, err)

	tasks, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Task 1", tasks[0].Title)
	assert.Equal(t, "Task 2", tasks[1].Title)
}

func TestTaskRepository_GetAll_Empty(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	ctx := context.Background()

	tasks, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tasks)
	assert.Len(t, tasks, 0)
}

func TestTaskRepository_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	ctx := context.Background()

	task := models.Task{
		Title:       "Original Title",
		Description: "Original Description",
		Status:      models.New,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, task)
	require.NoError(t, err)

	updatedTask := models.Task{
		Title:       "Updated Title",
		Description: "Updated Description",
		Status:      models.InProgress,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.Update(ctx, 1, updatedTask)
	assert.NoError(t, err)

	retrievedTask, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", retrievedTask.Title)
	assert.Equal(t, "Updated Description", retrievedTask.Description)
	assert.Equal(t, models.InProgress, retrievedTask.Status)
}

func TestTaskRepository_Update_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	ctx := context.Background()

	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.New,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Update(ctx, 999, task)
	assert.Error(t, err)
}

func TestTaskRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	ctx := context.Background()

	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.New,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, task)
	require.NoError(t, err)

	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM tasks").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestTaskRepository_Delete_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	assert.Error(t, err)
}
