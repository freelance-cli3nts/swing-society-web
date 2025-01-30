package routes

import (
    "net/http"
    
    "your-username/swing-society/internal/handlers"
    "your-username/swing-society/internal/static"
)

// Setup configures all routes
func Setup() {
    // Setup static file serving
    static.ServeFiles()

    // Handle all template requests
    http.HandleFunc("/templates/", handlers.HandleTemplate)
}