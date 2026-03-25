package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	v2 "math/rand/v2"

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
func (s *Service) SaveFile(fileData io.Reader, filename string) (string, error) {
	data, err := io.ReadAll(fileData)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	var hash string

	// FIXME: maybe there is a better way
	for {
		hash = randomString(3)
		_, err := s.repository.HeadFile(hash)
		if err != nil {
			break
		}
	}

	err = s.repository.UploadFile(hash, filename, data)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return hash, nil
}

var ra *v2.ChaCha8

func init() {
	buf := [32]byte{}
	rand.Read(buf[:])
	ra = v2.NewChaCha8(buf)
}

func randomString(len int) string {
	buf := make([]byte, len)
	ra.Read(buf)
	return hex.EncodeToString(buf)
}
