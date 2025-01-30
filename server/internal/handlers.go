package internal

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
)

// HandleTemplate handles all template requests
func HandleTemplate(w http.ResponseWriter, r *http.Request) {
    log.Printf("\n=== New Request ===")
    log.Printf("Request URL: %s", r.URL.Path)
    
    // Get the root directory path
    rootDir, err := filepath.Abs("../")
    if err != nil {
        log.Printf("Error getting root directory: %v", err)
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    
    // Construct full file path
    filePath := filepath.Join(rootDir, r.URL.Path)
    log.Printf("Trying to serve file: %s", filePath)
    
    // Check if file exists
    _, err = os.Stat(filePath)
    if os.IsNotExist(err) {
        log.Printf("File not found: %s", filePath)
        http.NotFound(w, r)
        return
    }
    
    // Set headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "text/html")
    
    // Serve the file
    http.ServeFile(w, r, filePath)
    log.Printf("File served successfully")
}