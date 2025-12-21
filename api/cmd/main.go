package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moabdelazem/mutlitier_app/internal/config"
	"github.com/moabdelazem/mutlitier_app/internal/database"
	"github.com/moabdelazem/mutlitier_app/internal/handler"
	"github.com/moabdelazem/mutlitier_app/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize structured logger
	log := logger.Init(&cfg.LogConfig)

	log.Info().
		Str("environment", cfg.Environment).
		Str("port", cfg.SrvPort).
		Str("log_level", cfg.LogConfig.Level).
		Str("log_format", cfg.LogConfig.Format).
		Msg("Starting application")

	// Connect to database
	db, err := database.NewPostgresConnection(&cfg.DatabaseConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	log.Info().Msg("Database connection established")

	// Setup router with config and logger
	router := handler.SetupRouter(db, cfg, log)

	// Configure HTTP server
	srv := &http.Server{
		Addr:           cfg.SrvPort,
		Handler:        router,
		ReadTimeout:    time.Second * 15,
		WriteTimeout:   time.Second * 15,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 1 << 20, // 1mb
	}

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info().Str("addr", cfg.SrvPort).Msg("Server started")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	<-quit
	log.Info().Msg("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	if err := db.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing database")
	}

	log.Info().Msg("Server stopped")
}
