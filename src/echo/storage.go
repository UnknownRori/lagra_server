package echo

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type StorageService struct {
	path string
}

func NewStorageService(path string) *StorageService {
	s := StorageService{path: path}
	return &s
}

func (s *StorageService) store(file *multipart.FileHeader) (string, error) {
	uuid := uuid.New()
	filename := fmt.Sprintf("%s/%s%s", s.path, uuid, filepath.Ext(file.Filename))
	src, err := file.Open()
	if err != nil {
		return filename, err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(filename)
	if err != nil {
		return filename, err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return filename, err
	}

	return filename, nil
}
