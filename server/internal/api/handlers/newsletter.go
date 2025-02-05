package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "swinng-society-website/server/internal/api/models"
)

func HandleNewsletterSubscription(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var subscription models.Newsletter
    if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate required fields
    if subscription.Email == "" {
        http.Error(w, "Email is required", http.StatusBadRequest)
        return
    }

    // TODO: Add email validation
    // TODO: Add actual newsletter subscription logic
    // TODO: Check for duplicate subscriptions

    log.Printf("Newsletter subscription received: %+v", subscription)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Thanks for subscribing to our newsletter!",
    })
}