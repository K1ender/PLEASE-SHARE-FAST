package logger

import (
	"context"
	"log/slog"
	"net/http"
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

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.Default().With(
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"request_id", r.Header.Get("X-Request-Id"),
		)

		ctx := context.WithValue(r.Context(), loggerKey{}, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
