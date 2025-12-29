package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// responseWriter wrapper untuk menangkap status code dan bytes
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// Logger middleware untuk logging HTTP requests
func Logger(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Skip logging untuk health check dan static files
			if r.URL.Path == "/health" || r.URL.Path == "/ping" {
				next.ServeHTTP(w, r)
				return
			}

			// Wrap response writer
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Proses request
			next.ServeHTTP(wrapped, r)

			// Hitung durasi
			duration := time.Since(start)

			// Log dengan fields
			fields := logrus.Fields{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      wrapped.statusCode,
				"duration_ms": duration.Milliseconds(),
				"ip":          r.RemoteAddr,
				"user_agent":  r.UserAgent(),
			}

			// Tambahkan query params jika ada
			if r.URL.RawQuery != "" {
				fields["query"] = r.URL.RawQuery
			}

			// Tentukan level log berdasarkan status code
			entry := logger.WithFields(fields)
			switch {
			case wrapped.statusCode >= 500:
				entry.Error("Server Error")
			case wrapped.statusCode >= 400:
				entry.Warn("Client Error")
			case wrapped.statusCode >= 300:
				entry.Info("Redirect")
			default:
				entry.Info("Request processed")
			}
		})
	}
}

// Recovery middleware untuk menangkap panic
func Recovery(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.WithFields(logrus.Fields{
						"error":  err,
						"method": r.Method,
						"path":   r.URL.Path,
						"ip":     r.RemoteAddr,
					}).Error("Panic recovered")

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"Internal Server Error"}`))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
