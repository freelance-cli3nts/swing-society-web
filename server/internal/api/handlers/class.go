package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "swinng-society-website/server/internal/api/models"
)

func HandleClassInquiry(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var inquiry models.ClassInquiry
    if err := json.NewDecoder(r.Body).Decode(&inquiry); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate required fields
    if inquiry.Name == "" || inquiry.Email == "" || inquiry.ClassType == "" || inquiry.Level == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // TODO: Add email validation
    // TODO: Add actual class registration logic

    log.Printf("Class inquiry received: %+v", inquiry)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Thank you for your interest! We'll contact you with class details soon.",
    })
}