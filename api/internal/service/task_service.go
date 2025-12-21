package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/moabdelazem/mutlitier_app/internal/model"
	"github.com/moabdelazem/mutlitier_app/internal/repository"
)

var (
	ErrValidation   = errors.New("validation error")
	ErrTaskNotFound = errors.New("task not found")
)

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// TaskService handles business logic for tasks
type TaskService struct {
	repo     *repository.TaskRepository
	validate *validator.Validate
}

// NewTaskService creates a new TaskService
func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{
		repo:     repo,
		validate: validator.New(),
	}
}

// Create creates a new task
func (s *TaskService) Create(ctx context.Context, req *model.CreateTaskRequest) (*model.TaskResponse, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, formatValidationErrors(err))
	}

	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
	}

	createdTask, err := s.repo.Create(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return createdTask.ToResponse(), nil
}

// GetByID retrieves a task by its ID
func (s *TaskService) GetByID(ctx context.Context, id string) (*model.TaskResponse, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task.ToResponse(), nil
}

// GetAll retrieves all tasks
func (s *TaskService) GetAll(ctx context.Context) ([]*model.TaskResponse, error) {
	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var responses []*model.TaskResponse
	for _, task := range tasks {
		responses = append(responses, task.ToResponse())
	}

	// Return empty slice instead of nil
	if responses == nil {
		responses = []*model.TaskResponse{}
	}

	return responses, nil
}

// Update updates a task
func (s *TaskService) Update(ctx context.Context, id string, req *model.UpdateTaskRequest) (*model.TaskResponse, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, formatValidationErrors(err))
	}

	updatedTask, err := s.repo.Update(ctx, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return updatedTask.ToResponse(), nil
}

// Delete deletes a task
func (s *TaskService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// formatValidationErrors formats validation errors into a user-friendly message
func formatValidationErrors(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var messages []string
		for _, e := range validationErrors {
			var message string
			switch e.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", e.Field())
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
			case "oneof":
				message = fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
			default:
				message = fmt.Sprintf("%s is invalid", e.Field())
			}
			messages = append(messages, message)
		}
		if len(messages) > 0 {
			return messages[0] // Return first error for simplicity
		}
	}
	return err.Error()
}
