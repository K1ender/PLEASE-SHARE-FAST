package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k1ender/psf/internal/cleaner"
	"github.com/k1ender/psf/internal/config"
	"github.com/k1ender/psf/internal/logger"
	"github.com/k1ender/psf/internal/repository"
	"github.com/k1ender/psf/internal/service"
	httptransport "github.com/k1ender/psf/internal/transport/http"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
)

func Run(ctx context.Context) error {
	cfg := config.MustInit()

	zaplog := zap.Must(zap.NewProduction())
	log := slog.New(slogzap.Option{Level: slog.LevelInfo, Logger: zaplog}.NewZapHandler())

	ctx = logger.WithLogger(ctx, log)

	log.Warn("You using in-memory repository, it will delete all files after restart", slog.String("repository", "in-memory"))

	filerepo := repository.NewInMemoryRepository()
	fileService := service.NewService(filerepo)

	clean := cleaner.NewInMemoryCleaner(fileService)

	http := httptransport.New(cfg.HTTP, fileService)

	log.Info("starting server")

	ticker := time.NewTicker(time.Hour)

	go func() {
		for range ticker.C {
			clean.Clean(ctx)
		}
	}()

	go func() {
		err := http.Run(ctx)
		if err != nil {
			log.Error("failed to start server", slog.Any("error", err))
			panic(err)
		}
	}()

	log.Info("server started", slog.String("address", ":8080"))

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	log.Info("shutting down server")

	err := http.Shutdown(ctx)
	if err != nil {
		log.Error("failed to shutdown server", slog.Any("error", err))
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	log.Info("server shutdown")

	return nil
}
