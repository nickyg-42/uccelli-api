package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware handles request logging with detailed information
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process the request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Get username from context if available
		username := "unknown"
		if user, ok := r.Context().Value("username").(string); ok {
			username = user
		}

		// Log the request details
		log.Printf(
			"[%s] %s %s %s | Status: %d | Duration: %v | IP: %s | User: %s | User-Agent: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			r.Proto,
			rw.statusCode,
			duration,
			getClientIP(r),
			username,
			r.UserAgent(),
		)
	})
}

// responseWriter is a custom ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}
