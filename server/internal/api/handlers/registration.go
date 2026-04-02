// internal/api/handlers/registration.go
package handlers

import (
    "log"
    "net/http"
    "swing-society-website/server/internal/api/models"
    "swing-society-website/server/internal/api/response"
    "swing-society-website/server/internal/email"
    "swing-society-website/server/internal/middleware"
    customerrors "swing-society-website/server/internal/errors"
    "swing-society-website/server/internal/storage"
)

type RegistrationHandler struct {
    storage  storage.RegistrationStorage
    emailSvc *email.Service
}

func NewRegistrationHandler(storage storage.RegistrationStorage, emailSvc *email.Service) *RegistrationHandler {
    return &RegistrationHandler{
        storage:  storage,
        emailSvc: emailSvc,
    }
}

func (h *RegistrationHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {

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
    form.Phone = r.FormValue("phone")
    form.Roles = r.Form["roles[]"]
    form.RegisterAlone = r.FormValue("registerAlone")
    form.PartnerName = r.FormValue("partner-name")
    form.Source = r.FormValue("source")
    form.Message = r.FormValue("message")
    form.Timestamp = middleware.GetRequestTime(r.Context())

    if validationErrors := form.Validate(); len(validationErrors) > 0 {
        response.Error(w, customerrors.NewValidationError(
            "Form validation failed",
            validationErrors,
        ))
        return
    }

    // Store registration
    if err := h.storage.StoreRegistration(&form); err != nil {
        log.Printf("Failed to store registration: %v", err)
        response.Error(w, customerrors.NewInternalError(
            "Failed to store registration",
            err,
        ))
        return
    }

    // Send emails asynchronously so they don't block the response
    go func() {
        if err := h.emailSvc.SendWelcome(form.Name, form.Email); err != nil {
            log.Printf("Failed to send welcome email to %s: %v", form.Email, err)
        }
        if err := h.emailSvc.SendRegistrationNotification(form.Name, form.Email, form.Phone); err != nil {
            log.Printf("Failed to send registration notification: %v", err)
        }
    }()

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


// ValidateEmail handles the email validation for real-time inline validation
func (h *RegistrationHandler) ValidateEmail(w http.ResponseWriter, r *http.Request) {
	var form models.RegistrationForm
	
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
	}
	

	// Handle form-encoded request (from HTMX)
	if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte("<span class='error'>Invalid form data</span>"))
		} else {
				response.Error(w, customerrors.NewValidationError(
						"Invalid form format",
						map[string]string{"format": "Invalid form data"},
				))
		}
		return
	}

	form.Email = r.FormValue("email")
	log.Printf("Processing email validation for: %s", form.Email)
	
	// Run validation
	validationErrors := form.Validate()
	
 // If there are validation errors for email, return the error
 if emailError, exists := validationErrors["email"]; exists {
		if r.Header.Get("HX-Request") == "true" {
				// For HTMX requests, return HTML error message
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte("<span class='error'>" + emailError + "</span>"))
		} else {
				// For API requests, return JSON error
				response.Error(w, customerrors.NewValidationError(
						"Please enter a valid email address",
						map[string]string{"email": emailError},
				))
		}
			return
	}
	
	// If no validation errors for email, send success response
	if r.Header.Get("HX-Request") == "true" {
			// If it's an HTMX request, return empty to clear the error message
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(""))
	} else {
			// Regular API response
			response.JSON(w, http.StatusOK, map[string]string{"message": "Valid email"})
	}
}

// ValidateName handles the name validation for real-time inline validation
func (h *RegistrationHandler) ValidateName(w http.ResponseWriter, r *http.Request) {
	var form models.RegistrationForm
	
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
	}

	// Handle form-encoded request (from HTMX)
	if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte("<span class='error'>Invalid form data</span>"))
			} else
			{
				response.Error(w, customerrors.NewValidationError(
						"Invalid form format",
						map[string]string{"format": "Invalid form data"},
				))
			}
			return
	}

	form.Name = r.FormValue("name")
	log.Printf("Processing name validation for: %s", form.Name)

	// Run validation
	validationErrors := form.Validate()

	// If there are validation errors for name, return the error
	if nameError, exists := validationErrors["name"]; exists {
			if r.Header.Get("HX-Request") == "true" {
					// For HTMX requests, return HTML error message
					w.Header().Set("Content-Type", "text/html")
					w.Write([]byte("<span class='error'>" + nameError + "</span>"))
			} else
			{
					// For API requests, return JSON error
					response.Error(w, customerrors.NewValidationError(
							"Please enter a valid name",
							map[string]string{"name": nameError},
					))
			}
			return
	}

	// If no validation errors for name, send success response
	if r.Header.Get("HX-Request") == "true" {
			// If it's an HTMX request, return empty to clear the error message
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(""))
	} else {
			// Regular API response
			response.JSON(w, http.StatusOK, map[string]string{"message": "Valid name"})
	}
}


// ValidatePhone handles the phone validation for real-time inline validation
func (h *RegistrationHandler) ValidatePhone(w http.ResponseWriter, r *http.Request) {
	var form models.RegistrationForm
	
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
	}
	
	// Handle form-encoded request (from HTMX)
	if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("Content-Type", "text/html")
					w.Write([]byte("<span class='error'>Invalid form data</span>"))
			} else {
					response.Error(w, customerrors.NewValidationError(
							"Invalid form format",
							map[string]string{"format": "Invalid form data"},
					))
			}
			return
	}
	
	form.Phone = r.FormValue("phone")
	log.Printf("Processing phone validation for: %s", form.Phone)
	
	// Run validation
	validationErrors := form.Validate()
	
	// If there are validation errors for phone, return the error
	if phoneError, exists := validationErrors["phone"]; exists {
			if r.Header.Get("HX-Request") == "true" {
					// For HTMX requests, return HTML error message
					w.Header().Set("Content-Type", "text/html")
					w.Write([]byte("<span class='error'>" + phoneError + "</span>"))
			} else {
					// For API requests, return JSON error
					response.Error(w, customerrors.NewValidationError(
							"Please enter a valid phone number",
							map[string]string{"phone": phoneError},
					))
			}
			return
	}
	
	// If no validation errors for phone, send success response
	if r.Header.Get("HX-Request") == "true" {
			// If it's an HTMX request, return empty to clear the error message
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(""))
	} else {
			// Regular API response
			response.JSON(w, http.StatusOK, map[string]string{"message": "Valid phone"})
	}
}


