package handlers

import (
    "encoding/json"
    "net/http"
    "swing-society-website/server/internal/api/models"
    "swing-society-website/server/internal/api/response"
    customerrors "swing-society-website/server/internal/errors"
		"swing-society-website/server/internal/storage"
)
type NewsletterHandler struct {
	storage storage.NewsletterStorage
}

func NewNewsletterHandler(storage storage.NewsletterStorage) *NewsletterHandler {
	return &NewsletterHandler{
			storage: storage,
	}
}

func (h *NewsletterHandler) HandleNewsletter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
			response.Error(w, customerrors.NewValidationError(
					"Method not allowed",
					map[string]string{"method": "Only POST method is allowed"},
			))
			return
	}

	var newsletter models.Newsletter
	if err := json.NewDecoder(r.Body).Decode(&newsletter); err != nil {
			response.Error(w, customerrors.NewValidationError(
					"Invalid request body",
					map[string]string{"body": "Could not parse request body"},
			))
			return
	}

	if errors := newsletter.Validate(); len(errors) > 0 {
			response.Error(w, customerrors.NewValidationError(
					"Validation failed",
					errors,
			))
			return
	}

	if err := h.storage.Subscribe(&newsletter); err != nil {
			response.Error(w, customerrors.NewInternalError(
					"Failed to store newsletter subscription",
					err,
			))
			return
	}

	response.JSON(w, http.StatusOK, map[string]string{
			"message": "Thanks for subscribing to our newsletter!",
	})
}