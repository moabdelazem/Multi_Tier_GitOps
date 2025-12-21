package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/moabdelazem/mutlitier_app/internal/config"
	"github.com/moabdelazem/mutlitier_app/internal/database"
	"github.com/moabdelazem/mutlitier_app/internal/repository"
	"github.com/moabdelazem/mutlitier_app/internal/service"
	"github.com/moabdelazem/mutlitier_app/pkg"
	"github.com/moabdelazem/mutlitier_app/pkg/logger"
	"github.com/moabdelazem/mutlitier_app/pkg/middleware"
)

type HealthResponse struct {
	Status   string                 `json:"status"`
	Services map[string]ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Status  string         `json:"status"`
	Message string         `json:"message,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

type HealthHandler struct {
	db *database.DB
}

func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func SetupRouter(db *database.DB, cfg *config.Config, log *logger.Logger) http.Handler {
	r := chi.NewRouter()

	// Initialize handlers
	healthHandler := NewHealthHandler(db)

	// Initialize task dependencies
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := NewTaskHandler(taskService)

	// Core middlewares
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(60 * time.Second))

	// CORS middleware (configured via environment)
	r.Use(middleware.CORS(&cfg.CORSConfig))

	// Structured request logging (replaces chi's DefaultLogger)
	r.Use(middleware.RequestLogger(log))

	// Health check route
	r.Get("/health", healthHandler.healthCheckHandler)

	// Task routes
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", taskHandler.Create)
		r.Get("/", taskHandler.GetAll)
		r.Get("/{id}", taskHandler.GetByID)
		r.Put("/{id}", taskHandler.Update)
		r.Delete("/{id}", taskHandler.Delete)
	})

	return r
}

func (h *HealthHandler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	healthResp := HealthResponse{
		Status:   "healthy",
		Services: make(map[string]ServiceInfo),
	}

	dbStatus := h.checkDatabase(ctx)
	healthResp.Services["database"] = dbStatus

	// If database is down, overall status is unhealthy
	if dbStatus.Status == "unhealthy" {
		healthResp.Status = "unhealthy"
		pkg.ServiceUnavailable(w, "Service unhealthy")
		return
	}

	pkg.JSONSuccess(w, healthResp)
}

func (h *HealthHandler) checkDatabase(ctx context.Context) ServiceInfo {
	if err := h.db.PingContext(ctx); err != nil {
		return ServiceInfo{
			Status:  "unhealthy",
			Message: "Failed to ping database",
			Details: map[string]any{
				"error": err.Error(),
			},
		}
	}

	stats := h.db.Stats()

	return ServiceInfo{
		Status:  "healthy",
		Message: "Database connection is active",
		Details: map[string]any{
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
			"max_open":         stats.MaxOpenConnections,
			"wait_count":       stats.WaitCount,
			"wait_duration":    stats.WaitDuration.String(),
		},
	}
}
