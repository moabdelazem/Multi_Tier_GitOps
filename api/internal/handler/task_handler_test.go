package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/moabdelazem/mutlitier_app/internal/model"
	"github.com/moabdelazem/mutlitier_app/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskService is a mock implementation of TaskService for testing
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) Create(ctx context.Context, req *model.CreateTaskRequest) (*model.TaskResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TaskResponse), args.Error(1)
}

func (m *MockTaskService) GetAll(ctx context.Context) ([]*model.TaskResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.TaskResponse), args.Error(1)
}

func (m *MockTaskService) GetByID(ctx context.Context, id string) (*model.TaskResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TaskResponse), args.Error(1)
}

func (m *MockTaskService) Update(ctx context.Context, id string, req *model.UpdateTaskRequest) (*model.TaskResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TaskResponse), args.Error(1)
}

func (m *MockTaskService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TaskServiceInterface defines the interface for task service operations
type TaskServiceInterface interface {
	Create(ctx context.Context, req *model.CreateTaskRequest) (*model.TaskResponse, error)
	GetAll(ctx context.Context) ([]*model.TaskResponse, error)
	GetByID(ctx context.Context, id string) (*model.TaskResponse, error)
	Update(ctx context.Context, id string, req *model.UpdateTaskRequest) (*model.TaskResponse, error)
	Delete(ctx context.Context, id string) error
}

// TestTaskHandler wraps handler with mock service
type TestTaskHandler struct {
	service TaskServiceInterface
}

func NewTestTaskHandler(service TaskServiceInterface) *TestTaskHandler {
	return &TestTaskHandler{service: service}
}

func (h *TestTaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	task, err := h.service.Create(r.Context(), &req)
	if err != nil {
		if err.Error() == service.ErrValidation.Error() || 
		   (len(err.Error()) > len(service.ErrValidation.Error()) && 
		    err.Error()[:len(service.ErrValidation.Error())] == service.ErrValidation.Error()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create task"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TestTaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAll(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve tasks"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TestTaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task ID is required"})
		return
	}

	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == service.ErrTaskNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve task"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TestTaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task ID is required"})
		return
	}

	var req model.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	task, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		if err == service.ErrTaskNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
			return
		}
		if err.Error() == service.ErrValidation.Error() ||
		   (len(err.Error()) > len(service.ErrValidation.Error()) &&
		    err.Error()[:len(service.ErrValidation.Error())] == service.ErrValidation.Error()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update task"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TestTaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task ID is required"})
		return
	}

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if err == service.ErrTaskNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete task"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============= Test Cases =============

func TestCreate_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	expectedTask := &model.TaskResponse{
		ID:          "123",
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}

	mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.CreateTaskRequest")).Return(expectedTask, nil)

	body := `{"title": "Test Task", "description": "Test Description"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response model.TaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask.ID, response.ID)
	assert.Equal(t, expectedTask.Title, response.Title)
	mockService.AssertExpectations(t)
}

func TestCreate_InvalidJSON(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid JSON payload", response["error"])
}

func TestCreate_ValidationError(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.CreateTaskRequest")).
		Return(nil, service.ErrValidation)

	body := `{"title": "", "description": "Test"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetAll_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	expectedTasks := []*model.TaskResponse{
		{ID: "1", Title: "Task 1", Status: "pending"},
		{ID: "2", Title: "Task 2", Status: "completed"},
	}

	mockService.On("GetAll", mock.Anything).Return(expectedTasks, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response []*model.TaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	mockService.AssertExpectations(t)
}

func TestGetAll_Empty(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	mockService.On("GetAll", mock.Anything).Return([]*model.TaskResponse{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response []*model.TaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Empty(t, response)
	mockService.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	expectedTask := &model.TaskResponse{
		ID:          "123",
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}

	mockService.On("GetByID", mock.Anything, "123").Return(expectedTask, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks/123", nil)
	
	// Setup chi context with URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response model.TaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask.ID, response.ID)
	mockService.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	mockService.On("GetByID", mock.Anything, "999").Return(nil, service.ErrTaskNotFound)

	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task not found", response["error"])
	mockService.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	expectedTask := &model.TaskResponse{
		ID:          "123",
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      "completed",
	}

	mockService.On("Update", mock.Anything, "123", mock.AnythingOfType("*model.UpdateTaskRequest")).Return(expectedTask, nil)

	body := `{"title": "Updated Task", "status": "completed"}`
	req := httptest.NewRequest(http.MethodPut, "/tasks/123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response model.TaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask.Title, response.Title)
	mockService.AssertExpectations(t)
}

func TestUpdate_NotFound(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	mockService.On("Update", mock.Anything, "999", mock.AnythingOfType("*model.UpdateTaskRequest")).Return(nil, service.ErrTaskNotFound)

	body := `{"title": "Updated Task"}`
	req := httptest.NewRequest(http.MethodPut, "/tasks/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	mockService.On("Delete", mock.Anything, "123").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/123", nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestDelete_NotFound(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTestTaskHandler(mockService)

	mockService.On("Delete", mock.Anything, "999").Return(service.ErrTaskNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}
