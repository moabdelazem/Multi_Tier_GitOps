package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moabdelazem/mutlitier_app/internal/model"
	"github.com/moabdelazem/mutlitier_app/internal/service"
	"github.com/moabdelazem/mutlitier_app/pkg"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	service *service.TaskService
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// Create handles POST /tasks
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.BadRequest(w, "Invalid JSON payload")
		return
	}

	task, err := h.service.Create(r.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			pkg.BadRequest(w, err.Error())
			return
		}
		pkg.InternalError(w, "Failed to create task")
		return
	}

	pkg.Created(w, task)
}

// GetAll handles GET /tasks
func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAll(r.Context())
	if err != nil {
		pkg.InternalError(w, "Failed to retrieve tasks")
		return
	}

	pkg.JSONSuccess(w, tasks)
}

// GetByID handles GET /tasks/{id}
func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		pkg.BadRequest(w, "Task ID is required")
		return
	}

	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			pkg.NotFound(w, "Task not found")
			return
		}
		pkg.InternalError(w, "Failed to retrieve task")
		return
	}

	pkg.JSONSuccess(w, task)
}

// Update handles PUT /tasks/{id}
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		pkg.BadRequest(w, "Task ID is required")
		return
	}

	var req model.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.BadRequest(w, "Invalid JSON payload")
		return
	}

	task, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			pkg.BadRequest(w, err.Error())
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			pkg.NotFound(w, "Task not found")
			return
		}
		pkg.InternalError(w, "Failed to update task")
		return
	}

	pkg.JSONSuccess(w, task)
}

// Delete handles DELETE /tasks/{id}
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		pkg.BadRequest(w, "Task ID is required")
		return
	}

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			pkg.NotFound(w, "Task not found")
			return
		}
		pkg.InternalError(w, "Failed to delete task")
		return
	}

	pkg.NoContent(w)
}
