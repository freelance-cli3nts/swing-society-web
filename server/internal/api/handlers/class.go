package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "swing-society-website/server/internal/api/models"
)

func HandleClass(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var class models.ClassInquiry
    if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if errors := class.Validate(); len(errors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errors)
        return
    }

    log.Printf("Class inquiry received: %+v", class)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Thank you for your interest! We'll contact you with class details soon.",
    })
}