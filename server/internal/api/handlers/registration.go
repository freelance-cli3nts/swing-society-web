package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "swing-society-website/server/internal/api/models"
)

func HandleRegistration(w http.ResponseWriter, r *http.Request) {
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
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form
    var form models.RegistrationForm
    if err := r.ParseForm(); err != nil {
        log.Printf("Error parsing form: %v", err)
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    form.Name = r.FormValue("name")
    form.Email = r.FormValue("email")

    if errors := form.Validate(); len(errors) > 0 {
        log.Printf("Form validation failed: %v", errors)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errors)
        return
    }

    log.Printf("Successful registration: %+v", form)

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(`
        <div class="success-message show" role="alert">
            <h3>Благодарим за регистрацията!</h3>
            <p>Ще се свържем с вас скоро.</p>
            <button class="close-btn" hx-get="" hx-target="closest .form-container" hx-trigger="click" hx-swap="delete">&times;</button>
        </div>
    `))
}
