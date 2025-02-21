// internal/api/handlers/registration.go
package handlers

import (
    "log"
    "net/http"
    "swing-society-website/server/internal/api/models"
    "swing-society-website/server/internal/api/response"
    customerrors "swing-society-website/server/internal/errors"
    "swing-society-website/server/internal/storage"
)

type RegistrationHandler struct {
    storage storage.RegistrationStorage
}

func NewRegistrationHandler(storage storage.RegistrationStorage) *RegistrationHandler {
    return &RegistrationHandler{
        storage: storage,
    }
}

func (h *RegistrationHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
    log.Printf("Registration form request received: %s", r.Method)

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
    var form models.RegistrationForm
    if err := r.ParseForm(); err != nil {
        response.Error(w, customerrors.NewValidationError(
            "Error parsing form",
            map[string]string{"form": "Invalid form data"},
        ))
        return
    }

    form.Name = r.FormValue("name")
    form.Email = r.FormValue("email")

    if validationErrors := form.Validate(); len(validationErrors) > 0 {
        response.Error(w, customerrors.NewValidationError(
            "Form validation failed",
            validationErrors,
        ))
        return
    }

    // Store registration
    if err := h.storage.StoreRegistration(&form); err != nil {
        response.Error(w, customerrors.NewInternalError(
            "Failed to store registration",
            err,
        ))
        return
    }

    // Handle HTMX request
    if r.Header.Get("HX-Request") == "true" {
        response.HTMLFragment(w, http.StatusOK, `
            <div class="success-message show" role="alert">
                <h3>Благодарим за регистрацията!</h3>
                <p>Ще се свържем с вас скоро.</p>
                <button class="close-btn" hx-get="" hx-target="closest .form-container" hx-trigger="click" hx-swap="delete">&times;</button>
            </div>
        `)
        return
    }

    // Regular API response
    response.JSON(w, http.StatusOK, map[string]string{
        "message": "Registration successful",
    })
}
