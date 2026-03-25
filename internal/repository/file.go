package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/k1ender/psf/internal/model"
)

var (
	ErrNotFound = fmt.Errorf("file not found")
)

type File interface {
	UploadFile(
		ctx context.Context,
		hash string,
		fileName string,
		data []byte,
	) error
	GetFile(ctx context.Context, id string) (model.File, error)
	HeadFile(ctx context.Context, id string) (model.FileMetadata, error)
	GetAllFiles(ctx context.Context) ([]string, error)
	DeleteFile(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	files  map[string]model.FileMetadata
	folder string
}

func NewInMemoryRepository() File {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		panic(err)
	}

	return &InMemoryRepository{
		files:  make(map[string]model.FileMetadata),
		folder: "data",
	}
}

// GetFile implements [File].
func (i *InMemoryRepository) GetFile(ctx context.Context, id string) (
	model.File,
	error,
) {
	f, ok := i.files[id]
	if !ok {
		return model.File{}, ErrNotFound
	}

	fileData, err := os.ReadFile(fmt.Sprintf("%s/%s", i.folder, id))
	if err != nil {
		return model.File{}, fmt.Errorf("failed to read file: %w", err)
	}
	if errors.Is(err, os.ErrNotExist) {
		return model.File{}, ErrNotFound
	}

	return model.File{
		Data:         fileData,
		FileMetadata: model.FileMetadata{Filename: f.Filename},
	}, nil
}

// UploadFile implements [File].
func (i *InMemoryRepository) UploadFile(
	ctx context.Context,
	hash string,
	fileName string,
	data []byte,
) error {
	err := os.WriteFile(fmt.Sprintf("%s/%s", i.folder, hash), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	i.files[hash] = model.FileMetadata{
		Filename:  fileName,
		CreatedAt: time.Now(),
	}
	return nil
}

// HeadFile implements [File].
func (i *InMemoryRepository) HeadFile(ctx context.Context, id string) (model.FileMetadata, error) {
	f, ok := i.files[id]
	if !ok {
		return model.FileMetadata{}, ErrNotFound
	}

	return f, nil
}

// GetAllFiles implements [File].
func (i *InMemoryRepository) GetAllFiles(ctx context.Context) ([]string, error) {
	files := make([]string, 0, len(i.files))
	for id := range i.files {
		files = append(files, id)
	}
	return files, nil
}

// DeleteFile implements [File].
func (i *InMemoryRepository) DeleteFile(ctx context.Context, id string) error {
	err := os.Remove(fmt.Sprintf("%s/%s", i.folder, id))
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	delete(i.files, id)

	return nil
}
