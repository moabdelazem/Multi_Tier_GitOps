package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/moabdelazem/mutlitier_app/internal/database"
	"github.com/moabdelazem/mutlitier_app/internal/model"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	db *database.DB
}

// NewTaskRepository creates a new TaskRepository
func NewTaskRepository(db *database.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create inserts a new task into the database
func (r *TaskRepository) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	query := `
		INSERT INTO tasks (title, description, status)
		VALUES ($1, $2, $3)
		RETURNING id, title, description, status, created_at, updated_at
	`

	var createdTask model.Task
	err := r.db.QueryRowContext(ctx, query,
		task.Title,
		task.Description,
		"pending",
	).Scan(
		&createdTask.ID,
		&createdTask.Title,
		&createdTask.Description,
		&createdTask.Status,
		&createdTask.CreatedAt,
		&createdTask.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &createdTask, nil
}

// GetByID retrieves a task by its ID
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var task model.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// GetAll retrieves all tasks from the database
func (r *TaskRepository) GetAll(ctx context.Context) ([]*model.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// Update updates a task in the database
func (r *TaskRepository) Update(ctx context.Context, id string, updates *model.UpdateTaskRequest) (*model.Task, error) {
	// First, get the current task
	currentTask, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if updates.Title != nil {
		currentTask.Title = *updates.Title
	}
	if updates.Description != nil {
		currentTask.Description = *updates.Description
	}
	if updates.Status != nil {
		currentTask.Status = *updates.Status
	}

	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, title, description, status, created_at, updated_at
	`

	var updatedTask model.Task
	err = r.db.QueryRowContext(ctx, query,
		currentTask.Title,
		currentTask.Description,
		currentTask.Status,
		time.Now(),
		id,
	).Scan(
		&updatedTask.ID,
		&updatedTask.Title,
		&updatedTask.Description,
		&updatedTask.Status,
		&updatedTask.CreatedAt,
		&updatedTask.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return &updatedTask, nil
}

// Delete removes a task from the database
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
