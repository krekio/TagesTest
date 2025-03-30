package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type FileStorage struct {
	storagePath string
	mu          sync.RWMutex
}

func NewFileStorage(storagePath string) (*FileStorage, error) {
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			return nil, fmt.Errorf("error creating storage directory: %s", err)
		}
	}

	return &FileStorage{storagePath: storagePath}, nil
}

func (s *FileStorage) Save(ctx context.Context, fileName string, content io.Reader) error {
	filepath := path.Join(s.storagePath, fileName)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}
	defer file.Close()
	_, err = io.Copy(file, content)
	if err != nil {
		return fmt.Errorf("error saving file: %s", err)
	}
	return nil
}

func (s *FileStorage) GetFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	filePath := filepath.Join(s.storagePath, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

func (s *FileStorage) ListFiles(ctx context.Context) ([]*API.FileInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	files, err := os.ReadDir(s.storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var result []*API.FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		result = append(result, &api.FileInfo{
			Filename:  file.Name(),
			CreatedAt: info.ModTime().Format(time.RFC3339),
			UpdatedAt: info.ModTime().Format(time.RFC3339),
		})
	}

	return result, nil
}

func (s *FileStorage) DeleteFile(ctx context.Context, filename string) error {
	filePath := filepath.Join(s.storagePath, filename)
	return os.Remove(filePath)
}
