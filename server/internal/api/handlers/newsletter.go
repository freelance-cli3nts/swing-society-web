package handlers

import (
    "log"
    "net/http"
    "time"
    "swing-society-website/server/internal/api/models"
    "swing-society-website/server/internal/api/response"
    "swing-society-website/server/internal/email"
    customerrors "swing-society-website/server/internal/errors"
    "swing-society-website/server/internal/storage"
)

type NewsletterHandler struct {
    storage  storage.NewsletterStorage
    emailSvc *email.Service
}

func NewNewsletterHandler(storage storage.NewsletterStorage, emailSvc *email.Service) *NewsletterHandler {
    return &NewsletterHandler{
        storage:  storage,
        emailSvc: emailSvc,
    }
}

func (h *NewsletterHandler) HandleNewsletter(w http.ResponseWriter, r *http.Request) {
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
	if err := r.ParseForm(); err != nil {
		response.Error(w, customerrors.NewValidationError(
			"Error parsing form",
			map[string]string{"form": "Invalid form data"},
		))
		return
	}

	// Create newsletter object
	newsletter := &models.Newsletter{
		Email:         r.FormValue("email"),
		Name:          r.FormValue("firstName"),
		Phone:         r.FormValue("phone"),
		Frequency:     r.FormValue("frequency"),
		TermsAccepted: r.FormValue("regulations") == "on",
		GdprAccepted:  r.FormValue("gdpr") == "on",
		Timestamp:     time.Now().Unix(),
	}

	// Validate the form
	if errors := newsletter.Validate(); len(errors) > 0 {
		response.Error(w, customerrors.NewValidationError(
			"Validation failed",
			errors,
		))
		return
	}

    if err := h.storage.Subscribe(newsletter); err != nil {
        response.Error(w, customerrors.NewInternalError(
            "Failed to store newsletter subscription",
            err,
        ))
        return
    }

    go func() {
        if err := h.emailSvc.SendNewsletterConfirmation(newsletter.Name, newsletter.Email); err != nil {
            log.Printf("Failed to send newsletter confirmation to %s: %v", newsletter.Email, err)
        }
    }()

    // Handle HTMX request
	if r.Header.Get("HX-Request") == "true" {
		response.HTMLFragment(w, http.StatusOK, `
			<div class="success-message show" role="alert">
				<h3>Благодарим за абонамента!</h3>
				<p>Вече сте абонирани за нашия бюлетин и ще получавате известия за предстоящи събития.</p>
				<button class="close-btn" hx-get="" hx-target="closest .form-container" hx-trigger="click" hx-swap="delete">&times;</button>
			</div>
		`)
		return
	}

	// Regular API response
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Thanks for subscribing to our newsletter!",
	})
}