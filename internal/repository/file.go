package repository

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	v2 "math/rand/v2"

	"github.com/k1ender/psf/internal/model"
)

type File interface {
	UploadFile(file model.File) (string, error)
	GetFile(id string) (model.File, error)
}

type InMemoryRepository struct {
	files map[string]model.File
}

func NewInMemoryRepository() File {
	return &InMemoryRepository{
		files: make(map[string]model.File),
	}
}

// GetFile implements [File].
func (i *InMemoryRepository) GetFile(id string) (
	model.File,
	error,
) {
	f, ok := i.files[id]
	if !ok {
		return model.File{}, fmt.Errorf("file not found")
	}
	return f, nil
}

// UploadFile implements [File].
func (i *InMemoryRepository) UploadFile(
	file model.File,
) (string, error) {
	hash := randomString(3)
	i.files[hash] = file
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
