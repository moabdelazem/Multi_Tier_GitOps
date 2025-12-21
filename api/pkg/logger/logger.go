package logger

import (
	"io"
	"os"
	"time"

	"github.com/moabdelazem/mutlitier_app/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger is a wrapper around zerolog.Logger
type Logger struct {
	zerolog.Logger
}

// Global logger instance
var globalLogger *Logger

// Init initializes the global logger with the given configuration
func Init(cfg *config.LogConfig) *Logger {
	// Set log level
	level := parseLogLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	// Set time format
	zerolog.TimeFieldFormat = getTimeFormat(cfg.TimeFormat)

	// Configure output writer based on format
	var writer io.Writer
	if cfg.Format == "console" {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
	} else {
		// Default to JSON format (ideal for Kubernetes/log aggregators)
		writer = os.Stdout
	}

	// Create logger with common fields
	logger := zerolog.New(writer).
		With().
		Timestamp().
		Caller().
		Logger()

	globalLogger = &Logger{Logger: logger}

	// Set as default logger
	log.Logger = logger

	return globalLogger
}

// Get returns the global logger instance
func Get() *Logger {
	if globalLogger == nil {
		// Initialize with defaults if not initialized
		globalLogger = &Logger{
			Logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
		}
	}
	return globalLogger
}

// WithRequestID creates a new logger with request ID context
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		Logger: l.With().Str("request_id", requestID).Logger(),
	}
}

// WithComponent creates a new logger with component context
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		Logger: l.With().Str("component", component).Logger(),
	}
}

// parseLogLevel converts string log level to zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}

// getTimeFormat returns the time format string
func getTimeFormat(format string) string {
	switch format {
	case "unix":
		return zerolog.TimeFormatUnix
	case "unixms":
		return zerolog.TimeFormatUnixMs
	case "unixmicro":
		return zerolog.TimeFormatUnixMicro
	case "rfc3339":
		return time.RFC3339
	case "rfc3339nano":
		return time.RFC3339Nano
	default:
		return time.RFC3339
	}
}
