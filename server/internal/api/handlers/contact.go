package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "swing-society-website/server/internal/api/models"
)

func HandleContact(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var contact models.ContactForm
    if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if errors := contact.Validate(); len(errors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errors)
        return
    }

    log.Printf("Contact form received: %+v", contact)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Thank you for your message. We'll get back to you soon!",
    })
}