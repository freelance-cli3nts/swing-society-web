package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "swing-society-website/server/internal/api/models"
    "swing-society-website/server/internal/api/response"
    "swing-society-website/server/internal/email"
    customerrors "swing-society-website/server/internal/errors"
    "swing-society-website/server/internal/storage"
)

type ContactHandler struct {
    storage  storage.ContactStorage
    emailSvc *email.Service
}

func NewContactHandler(storage storage.ContactStorage, emailSvc *email.Service) *ContactHandler {
    return &ContactHandler{
        storage:  storage,
        emailSvc: emailSvc,
    }
}

func (h *ContactHandler) HandleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
			response.Error(w, customerrors.NewValidationError(
					"Method not allowed",
					map[string]string{"method": "Only POST method is allowed"},
			))
			return
	}

	var contact models.ContactForm
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			response.Error(w, customerrors.NewValidationError(
					"Invalid request body",
					map[string]string{"body": "Could not parse request body"},
			))
			return
	}

	if errors := contact.Validate(); len(errors) > 0 {
			response.Error(w, customerrors.NewValidationError(
					"Validation failed",
					errors,
			))
			return
	}

    if err := h.storage.StoreContactForm(&contact); err != nil {
        response.Error(w, customerrors.NewInternalError(
            "Failed to store contact form",
            err,
        ))
        return
    }

    go func() {
        if err := h.emailSvc.SendContactNotification(contact.Name, contact.Email, contact.Message); err != nil {
            log.Printf("Failed to send contact notification: %v", err)
        }
    }()

    response.JSON(w, http.StatusOK, map[string]string{
        "message": "Thank you for your message. We'll get back to you soon!",
    })
}






