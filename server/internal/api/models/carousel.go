// internal/api/models/carousel.go
package models

// CarouselItem represents a single item in the carousel
type CarouselItem struct {
    ImageURL string `json:"image_url"`
    Caption  string `json:"caption"`
    Type     string `json:"type"`
}

// Validate performs validation on the carousel item
func (c *CarouselItem) Validate() map[string]string {
    errors := make(map[string]string)

    if c.ImageURL == "" {
        errors["image_url"] = "Image URL is required"
    }

    if c.Type == "" {
        errors["type"] = "Type is required"
    }

    return errors
}