// internal/storage/template.go
package storage

import (
    "os"
    "path/filepath"
    "fmt"
    "swing-society-website/server/internal/config"
)

// TemplateStorage defines the interface for template operations
type TemplateStorage interface {
    GetTemplate(path string) (string, error)
    TemplateExists(path string) bool
    GetIndexTemplate() (string, error)
}

type FileTemplateStorage struct {
    templatesDir string
}

func NewFileTemplateStorage() *FileTemplateStorage {
    return &FileTemplateStorage{
        templatesDir: config.AppConfig.Paths.TemplatesDir,
    }
}

func (s *FileTemplateStorage) GetTemplate(path string) (string, error) {
    fullPath := filepath.Join(s.templatesDir, path)
    content, err := os.ReadFile(fullPath)
    if err != nil {
        return "", fmt.Errorf("failed to read template: %v", err)
    }
    return string(content), nil
}

func (s *FileTemplateStorage) TemplateExists(path string) bool {
    fullPath := filepath.Join(s.templatesDir, path)
    _, err := os.Stat(fullPath)
    return err == nil
}

func (s *FileTemplateStorage) GetIndexTemplate() (string, error) {
    return s.GetTemplate("index.html")
}
