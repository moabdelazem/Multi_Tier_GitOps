package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/moabdelazem/mutlitier_app/pkg/logger"
)

// RequestLogger is a structured logging middleware using zerolog
func RequestLogger(log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Get or generate request ID
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = middleware.GetReqID(r.Context())
			}

			// Wrap response writer to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Add request ID to response header
			ww.Header().Set("X-Request-ID", requestID)

			// Process request
			next.ServeHTTP(ww, r)

			// Calculate duration
			duration := time.Since(start)

			// Log the request
			logEvent := log.Info()
			if ww.Status() >= 500 {
				logEvent = log.Error()
			} else if ww.Status() >= 400 {
				logEvent = log.Warn()
			}

			logEvent.
				Str("request_id", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Int("status", ww.Status()).
				Int("bytes_written", ww.BytesWritten()).
				Dur("duration", duration).
				Str("duration_human", duration.String()).
				Msg("HTTP request")
		})
	}
}
