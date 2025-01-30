package internal

import (
		"log"
    "net/http"
		"path/filepath"
)

// Setup configures all routes
func SetupRoutes() {
    // Setup static file serving
    ServeFiles()

// Handle root path and templates
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request path: %s", r.URL.Path)
        
        rootDir, _ := filepath.Abs("../")
        var filePath string

        switch r.URL.Path {
        case "/":
            // Serve index.html for root path
            filePath = filepath.Join(rootDir, "templates/index.html")
        case "/templates/":
            // Serve index.html for /templates/ path too
            filePath = filepath.Join(rootDir, "templates/index.html")
        default:
            // Handle other template requests
            filePath = filepath.Join(rootDir, r.URL.Path)
        }

        log.Printf("Attempting to serve: %s", filePath)
        http.ServeFile(w, r, filePath)
    })
}
