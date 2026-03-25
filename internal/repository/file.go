package repository

import (
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
		hash string,
		fileName string,
		data []byte,
	) error
	GetFile(id string) (model.File, error)
	HeadFile(id string) (model.FileMetadata, error)
}

type InMemoryRepository struct {
	files  map[string]model.FileMetadata
	folder string
}

func NewInMemoryRepository() File {
	return &InMemoryRepository{
		files:  make(map[string]model.FileMetadata),
		folder: "data",
	}
}

// GetFile implements [File].
func (i *InMemoryRepository) GetFile(id string) (
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
func (i *InMemoryRepository) HeadFile(id string) (model.FileMetadata, error) {
	f, ok := i.files[id]
	if !ok {
		return model.FileMetadata{}, ErrNotFound
	}

	return f, nil
}
