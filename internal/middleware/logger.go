package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/k1ender/psf/internal/logger"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	if w.StatusCode == 0 {
		w.StatusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()

		log := slog.Default().With(
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"request_id", RequestIDFromContext(ctx),
		)

		ctx = logger.WithLogger(ctx, log)

		rw := &CustomResponseWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		defer func() {
			duration := time.Since(start)

			if rec := recover(); rec != nil {
				log.ErrorContext(ctx, "panic recovered",
					slog.Any("error", rec),
					slog.Int("status", http.StatusInternalServerError),
					slog.Int64("duration_ms", duration.Milliseconds()),
				)

				http.Error(rw, "internal server error", http.StatusInternalServerError)
				return
			}

			log.InfoContext(ctx, "request completed",
				slog.Int("status", rw.StatusCode),
				slog.Int64("duration_ms", duration.Milliseconds()),
			)
		}()

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
