package telemetry

import (
	"log/slog"
	"net/http"
	"time"
)

// ResponseWriter wrapper to capture status code and size
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int64
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

// Middleware returns a robust telemetry middleware.
// It logs request lifecycle, duration, status, and size.
// Ideally, this would also emit to Prometheus.
func Middleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Panics are the enemy of reliability. Recover and log.
		defer func() {
			if rec := recover(); rec != nil {
				duration := time.Since(start)
				logger.Error("http_panic_recovered",
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
					"duration_ms", duration.Milliseconds(),
					"panic", rec,
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Structured logging for high-cardinality observability
		logger.Info("http_request_completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"size_bytes", rw.size,
			"duration_ms", duration.Milliseconds(),
			"user_agent", r.UserAgent(),
		)
	})
}
