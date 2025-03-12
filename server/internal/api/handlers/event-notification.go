package handlers

import (
	"log"
	"net/http"
	"swing-society-website/server/internal/api/models"
	"swing-society-website/server/internal/api/response"
	"swing-society-website/server/internal/middleware"
	customerrors "swing-society-website/server/internal/errors"
	"swing-society-website/server/internal/storage"
)

type EventNotificationHandler struct {
	storage storage.EventNotificationStorage
}

func NewEventNotificationHandler(storage storage.EventNotificationStorage) *EventNotificationHandler {
	return &EventNotificationHandler{
		storage: storage,
	}
}

func (h *EventNotificationHandler) HandleEventNotification(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		response.Error(w, customerrors.NewValidationError(
			"Method not allowed",
			map[string]string{"method": "Only POST method is allowed"},
		))
		return
	}

	// Parse form
	var notification models.EventNotification
	if err := r.ParseForm(); err != nil {
		response.Error(w, customerrors.NewValidationError(
			"Error parsing form",
			map[string]string{"form": "Invalid form data"},
		))
		return
	}

	// Extract form values
	notification.Email = r.FormValue("email")
	notification.FirstName = r.FormValue("firstName")
	notification.Phone = r.FormValue("phone")
	notification.ContactPreference = r.FormValue("contactPreference")
	notification.EventTypes = r.Form["eventTypes[]"]
	notification.Frequency = r.FormValue("frequency")
	notification.TermsAccepted = r.FormValue("regulations") == "on"
	notification.GdprAccepted = r.FormValue("gdpr") == "on"
	notification.Timestamp = middleware.GetRequestTime(r.Context())

	// Validate form
	if validationErrors := notification.Validate(); len(validationErrors) > 0 {
		response.Error(w, customerrors.NewValidationError(
			"Form validation failed",
			validationErrors,
		))
		return
	}

	// Store notification
	if err := h.storage.StoreEventNotification(&notification); err != nil {
		log.Printf("Failed to store event notification: %v", err)
		response.Error(w, customerrors.NewInternalError(
			"Failed to store event notification",
			err,
		))
		return
	}

	// Handle HTMX request
	if r.Header.Get("HX-Request") == "true" {
		response.HTMLFragment(w, http.StatusOK, `
			<div class="success-message show" role="alert">
				<h3>Благодарим за абонамента!</h3>
				<p>Вече сте абонирани за известия за събития и ще получавате информация според вашите предпочитания.</p>
				<button class="close-btn" hx-get="" hx-target="closest .form-container" hx-trigger="click" hx-swap="delete">&times;</button>
			</div>
		`)
		return
	}

	// Regular API response
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Event notification subscription successful",
	})
}