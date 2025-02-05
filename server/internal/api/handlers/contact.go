package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "swinng-society-website/server/internal/api/models"
)

func HandleContactSubmission(w http.ResponseWriter, r *http.Request) {
    // Only allow POST method
    if r.Method != http.MethodPost {
				log.Printf("Method not allowed: %s", r.Method)  // Add this line

        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse the request body
    var contact models.ContactForm
    if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate required fields
    if contact.Name == "" || contact.Email == "" || contact.Message == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // TODO: Add email validation
    // TODO: Add actual email sending logic

    // Log the submission (for now)
    log.Printf("Contact form submission received: %+v", contact)

    // Send success response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Thank you for your message. We'll get back to you soon!",
    })
}