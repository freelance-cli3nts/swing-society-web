package handlers

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
)

// HandleTemplate handles all template requests
func HandleTemplate(w http.ResponseWriter, r *http.Request) {
    // Print request details
    log.Printf("\n=== New Request ===")
    log.Printf("Request URL: %s", r.URL.Path)
    
    // Get absolute path to requested file
    absPath, err := filepath.Abs(filepath.Join(".", r.URL.Path))
    if err != nil {
        log.Printf("Error getting absolute path: %v", err)
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    log.Printf("Absolute path: %s", absPath)
    
    // Check if file exists
    _, err = os.Stat(absPath)
    if os.IsNotExist(err) {
        log.Printf("File does not exist at: %s", absPath)
        http.NotFound(w, r)
        return
    }
    
    // Set headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "text/html")
    
    // Serve the file
    http.ServeFile(w, r, absPath)
    log.Printf("File served successfully")
}