package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k1ender/psf/internal/repository"
	"github.com/k1ender/psf/internal/service"
	httptransport "github.com/k1ender/psf/internal/transport/http"
)

func Run(ctx context.Context) error {
	filerepo := repository.NewInMemoryRepository()
	fileService := service.NewService(filerepo)

	http := httptransport.New(":8080", fileService)

	slog.Info("starting server")

	ticker := time.NewTicker(time.Hour)

	go func() {
		for range ticker.C {
			files, err := fileService.GetAllFiles()
			if err != nil {
				slog.Error("failed to delete old files", slog.Any("error", err))
			}

			for _, id := range files {
				fileMetadata, err := fileService.HeadFile(id)
				if err != nil {
					slog.Error("failed to delete old files", slog.Any("error", err))
				}
				if time.Since(fileMetadata.CreatedAt) > 24*time.Hour {
					continue
				}

				err = fileService.DeleteFile(id)
				if err != nil {
					slog.Error("failed to delete old files", slog.Any("error", err))
				}
			}
		}
	}()

	go func() {
		err := http.Run(ctx)
		if err != nil {
			slog.Error("failed to start server", slog.Any("error", err))
			panic(err)
		}
	}()

	slog.Info("server started", slog.String("address", ":8080"))

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	slog.Info("shutting down server")

	err := http.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.Any("error", err))
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	slog.Info("server shutdown")

	return nil
}
