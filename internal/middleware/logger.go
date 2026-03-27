package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	slog.SetDefault(logger)
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok || logger == nil {
		return slog.Default()
	}

	return logger
}

func L(ctx context.Context) *slog.Logger {
	return FromContext(ctx)
}

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
		start := time.Now()

		ctx := r.Context()

		logger := slog.Default().With(
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"request_id", RequestIDFromContext(ctx),
		)

		ctx = context.WithValue(ctx, loggerKey{}, logger)

		rw := &CustomResponseWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		defer func() {
			duration := time.Since(start)

			if rec := recover(); rec != nil {
				logger.ErrorContext(ctx, "panic recovered",
					slog.Any("error", rec),
					slog.Int("status", http.StatusInternalServerError),
					slog.Int64("duration_ms", duration.Milliseconds()),
				)

				http.Error(rw, "internal server error", http.StatusInternalServerError)
				return
			}

			logger.InfoContext(ctx, "request completed",
				slog.Int("status", rw.StatusCode),
				slog.Int64("duration_ms", duration.Milliseconds()),
			)
		}()

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
