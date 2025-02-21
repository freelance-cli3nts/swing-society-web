// internal/api/handlers/carousel.go
package handlers

import (
    "log"
    "net/http"
    "strings"
		"swing-society-website/server/internal/api/models"
    "swing-society-website/server/internal/api/response"
    customerrors "swing-society-website/server/internal/errors"
    "swing-society-website/server/internal/storage"
)

// CarouselHandler handles carousel-related requests
type CarouselHandler struct {
    primaryStorage   storage.CarouselStorage
    fallbackStorage storage.CarouselStorage
}

// NewCarouselHandler creates a new CarouselHandler instance
func NewCarouselHandler(primary, fallback storage.CarouselStorage) *CarouselHandler {
    return &CarouselHandler{
        primaryStorage:   primary,
        fallbackStorage: fallback,
    }
}

// ServeCarousel handles the carousel data request
func (h *CarouselHandler) ServeCarousel(w http.ResponseWriter, r *http.Request) {
    carouselType := strings.TrimPrefix(r.URL.Path, "/api/carousel/")
    carouselType = strings.ToLower(carouselType)

    log.Printf("Received request for carousel data - Type: %s", carouselType)

    data, err := h.fetchCarouselData()
    if err != nil {
        response.Error(w, customerrors.NewInternalError(
            "Error fetching carousel data",
            err,
        ))
        return
    }

    filteredItems := make([]models.CarouselItem, 0)
    for _, item := range data {
        if item.Type == carouselType {
            filteredItems = append(filteredItems, item)
        }
    }

    response.JSON(w, http.StatusOK, filteredItems)
}

// fetchCarouselData attempts to fetch data from primary storage,
// falling back to secondary storage if needed
func (h *CarouselHandler) fetchCarouselData() ([]models.CarouselItem, error) {
    // Try primary storage (Google Sheets)
    data, err := h.primaryStorage.FetchItems()
    if err == nil {
        return data, nil
    }
    log.Printf("⚠️ Primary storage failed: %v", err)

    // Try fallback storage (JSON file)
    data, err = h.fallbackStorage.FetchItems()
    if err != nil {
        return nil, customerrors.NewInternalError(
            "Failed to fetch data from all storage sources",
            err,
        )
    }

    return data, nil
}