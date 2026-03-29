package cleaner

import (
	"context"
	"log/slog"
	"time"

	"github.com/k1ender/psf/internal/logger"
	"github.com/k1ender/psf/internal/service"
)

type Cleaner interface {
	Clean(ctx context.Context) error
}

type InMemoryCleaner struct {
	fileService service.File
}

func NewInMemoryCleaner(fileService service.File) Cleaner {
	return &InMemoryCleaner{fileService: fileService}
}

func (c *InMemoryCleaner) Clean(ctx context.Context) error {
	log := logger.FromContext(ctx)

	log.Info("cleaning up old files")

	files, err := c.fileService.GetAllFiles(ctx)
	if err != nil {
		log.Error("failed to delete old files", slog.Any("error", err))
	}

	for _, id := range files {
		fileMetadata, err := c.fileService.HeadFile(ctx, id)
		if err != nil {
			log.Error("failed to get file metadata", slog.Any("error", err))
			continue
		}
		if time.Since(fileMetadata.CreatedAt) <= 24*time.Hour {
			continue
		}

		err = c.fileService.DeleteFile(ctx, id)
		if err != nil {
			log.Error("failed to delete old file", slog.Any("error", err))
		}
	}

	return nil
}
