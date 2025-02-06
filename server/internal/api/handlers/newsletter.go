package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "swing-society-website/server/internal/api/models"
)

func HandleNewsletter(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var newsletter models.Newsletter
    if err := json.NewDecoder(r.Body).Decode(&newsletter); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if errors := newsletter.Validate(); len(errors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errors)
        return
    }

    log.Printf("Newsletter subscription received: %+v", newsletter)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Thanks for subscribing to our newsletter!",
    })
}