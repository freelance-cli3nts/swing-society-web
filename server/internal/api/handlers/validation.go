package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"swing-society-website/server/internal/api/response"
	customerrors "swing-society-website/server/internal/errors"
)

// ValidateEmail handles email validation for all forms
func ValidateEmail(w http.ResponseWriter, r *http.Request) {
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

	email := r.FormValue("email")
	log.Printf("Processing email validation for: %s", email)
	
	// Run validation
	var errorMsg string
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if strings.TrimSpace(email) == "" {
		errorMsg = "Моля, въведете вашия имейл"
	} else if !emailRegex.MatchString(strings.ToLower(email)) {
		errorMsg = "Моля, въведете валиден имейл адрес"
	}
	
	// If there are validation errors for email, return the error
	if errorMsg != "" {
		if r.Header.Get("HX-Request") == "true" {
			// For HTMX requests, return HTML error message
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<span class='error'>" + errorMsg + "</span>"))
		} else {
			// For API requests, return JSON error
			response.Error(w, customerrors.NewValidationError(
				"Моля, въведете валиден имейл адрес",
				map[string]string{"email": errorMsg},
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

// ValidatePhone handles phone validation for all forms
func ValidatePhone(w http.ResponseWriter, r *http.Request) {
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

	phone := r.FormValue("phone")
	log.Printf("Processing phone validation for: %s", phone)
	
	// Run validation
	var errorMsg string
	phoneRegex := regexp.MustCompile(`^(\+359|0|00359)[0-9]{9}$`)
	if strings.TrimSpace(phone) != "" && !phoneRegex.MatchString(strings.ToLower(phone)) {
		errorMsg = "Моля, въведете валиден телефонен номер: 0888123456 / +359888123456 / 00359888123456"
	}
	
	// If there are validation errors for phone, return the error
	if errorMsg != "" {
		if r.Header.Get("HX-Request") == "true" {
			// For HTMX requests, return HTML error message
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<span class='error'>" + errorMsg + "</span>"))
		} else {
			// For API requests, return JSON error
			response.Error(w, customerrors.NewValidationError(
				"Моля, въведете валиден БГ телефонен номер",
				map[string]string{"phone": errorMsg},
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

// ValidateName handles name validation for all forms
func ValidateName(w http.ResponseWriter, r *http.Request) {
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

	name := r.FormValue("name")
	log.Printf("Processing name validation for: %s", name)
	
	// Run validation
	var errorMsg string
	if strings.TrimSpace(name) == "" {
		errorMsg = "Моля, въведете вашето име"
	} else if len(name) < 3 {
		errorMsg = "Името трябва да е поне 3 символа"
	}
	
	// If there are validation errors for name, return the error
	if errorMsg != "" {
		if r.Header.Get("HX-Request") == "true" {
			// For HTMX requests, return HTML error message
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<span class='error'>" + errorMsg + "</span>"))
		} else {
			// For API requests, return JSON error
			response.Error(w, customerrors.NewValidationError(
				"Името трябва да е поне 3 символа",
				map[string]string{"name": errorMsg},
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