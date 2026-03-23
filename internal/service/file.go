package service

import (
	"fmt"
	"io"

	"github.com/k1ender/psf/internal/model"
	"github.com/k1ender/psf/internal/repository"
)

type File interface {
	SaveFile(file io.Reader, filename string) (string, error)
	GetFile(id string) ([]byte, string, error)
}

type Service struct {
	repository repository.File
}

func NewService(repository repository.File) File {
	return &Service{repository: repository}
}

// GetFile implements [File].
func (s *Service) GetFile(id string) ([]byte, string, error) {
	file, err := s.repository.GetFile(id)
	if err != nil {
		return nil, "", err
	}

	return file.Data, file.Filename, nil
}

// SaveFile implements [File].
func (s *Service) SaveFile(file io.Reader, filename string) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	files := model.File{
		Data:     data,
		Filename: filename,
	}

	hash, err := s.repository.UploadFile(files)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return hash, nil
}
