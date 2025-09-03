package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task models.Task) error {
	query := `
		INSERT INTO tasks (title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var id int
	err := r.db.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1`

	var task models.Task
	err := r.db.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, id int, task models.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5`

	result, err := r.db.Exec(ctx, query, task.Title, task.Description, task.Status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
