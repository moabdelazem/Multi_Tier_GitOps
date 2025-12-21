package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// App Configurations
type Config struct {
	SrvPort        string
	Environment    string
	DatabaseConfig DatabaseConfig
	CORSConfig     CORSConfig
	LogConfig      LogConfig
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// CORSConfig holds CORS settings - all configurable via environment variables
type CORSConfig struct {
	AllowedOrigins   []string // * or list of origins
	AllowedMethods   []string // GET, POST, PUT, DELETE, OPTIONS
	AllowedHeaders   []string // Accept, Authorization, Content-Type, X-Request-ID
	ExposedHeaders   []string // X-Request-ID
	AllowCredentials bool     
	MaxAge           int     
}

// LogConfig holds logging settings - configurable for different environments
type LogConfig struct {
	Level      string // LOG_LEVEL: debug, info, warn, error
	Format     string // LOG_FORMAT: json, console
	TimeFormat string // LOG_TIME_FORMAT: unix, rfc3339, etc.
}

// Create new config struct
func NewConfig() *Config {
	// Only load .env in development - in Kubernetes, env vars come from ConfigMap/Secrets
	if os.Getenv("ENVIRONMENT") != "production" {
		godotenv.Load()
	}

	return &Config{
		SrvPort:     getEnv("PORT", ":8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseConfig: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvAsInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "multi_tier_db"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),
		},
		CORSConfig: CORSConfig{
			AllowedOrigins:   getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods:   getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders:   getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"}),
			ExposedHeaders:   getEnvAsSlice("CORS_EXPOSED_HEADERS", []string{"X-Request-ID"}),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", false),
			MaxAge:           getEnvAsInt("CORS_MAX_AGE", 300),
		},
		LogConfig: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			TimeFormat: getEnv("LOG_TIME_FORMAT", "rfc3339"),
		},
	}
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// Get The Environment Variables
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
