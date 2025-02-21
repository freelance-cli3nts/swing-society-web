// internal/storage/carousel.go
package storage

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "swing-society-website/server/internal/api/models"
)

// CarouselStorage defines the interface for carousel data storage
type CarouselStorage interface {
    FetchItems() ([]models.CarouselItem, error)
}

// GoogleSheetsStorage implements CarouselStorage using Google Sheets
type GoogleSheetsStorage struct {
    SheetURL string
    Client   *http.Client
}

// NewGoogleSheetsStorage creates a new GoogleSheetsStorage instance
func NewGoogleSheetsStorage(sheetURL string) *GoogleSheetsStorage {
    return &GoogleSheetsStorage{
        SheetURL: sheetURL,
        Client:   &http.Client{},
    }
}

// FetchItems implements CarouselStorage interface for Google Sheets
func (s *GoogleSheetsStorage) FetchItems() ([]models.CarouselItem, error) {
    if s.SheetURL == "" {
        return nil, fmt.Errorf("Google Sheets URL is not set")
    }

    resp, err := s.Client.Get(s.SheetURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch from Google Sheets: %v", err)
    }
    defer resp.Body.Close()

    reader := csv.NewReader(resp.Body)
    var items []models.CarouselItem

    // Skip header row
    _, _ = reader.Read()

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("error reading CSV record: %v", err)
        }

        item := models.CarouselItem{
            ImageURL: convertGoogleDriveLink(record[0]),
            Caption:  record[1],
            Type:     strings.ToLower(record[2]),
        }
        
        if errors := item.Validate(); len(errors) == 0 {
            items = append(items, item)
        }
    }
    return items, nil
}

// JSONFileStorage implements CarouselStorage using local JSON file
type JSONFileStorage struct {
    FilePath string
}

// NewJSONFileStorage creates a new JSONFileStorage instance
func NewJSONFileStorage(filePath string) *JSONFileStorage {
    return &JSONFileStorage{
        FilePath: filePath,
    }
}

// FetchItems implements CarouselStorage interface for JSON file
func (s *JSONFileStorage) FetchItems() ([]models.CarouselItem, error) {
    file, err := os.Open(s.FilePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open JSON file: %v", err)
    }
    defer file.Close()

    var items []models.CarouselItem
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&items); err != nil {
        return nil, fmt.Errorf("failed to decode JSON: %v", err)
    }

    // Validate and filter items
    validItems := make([]models.CarouselItem, 0)
    for _, item := range items {
        if errors := item.Validate(); len(errors) == 0 {
            validItems = append(validItems, item)
        }
    }

    return validItems, nil
}

// Helper function to convert Google Drive links
func convertGoogleDriveLink(url string) string {
    if strings.Contains(url, "drive.google.com/file/d/") {
        parts := strings.Split(url, "/")
        if len(parts) >= 5 {
            fileID := parts[5]
            return "https://drive.google.com/uc?export=view&id=" + fileID
        }
    }
    return url
}