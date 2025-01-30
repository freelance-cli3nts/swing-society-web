package internal

import (
	"net/http"
	"path/filepath"
	"log"
	"os"
)

// ServeFiles sets up static file serving
// func ServeFiles() {
// 		rootDir, _ := filepath.Abs("../")
// 		staticDir := http.Dir(filepath.Join(rootDir,"static"))
//     fs := http.FileServer(staticDir))
//     http.Handle("/static/", http.StripPrefix("/static/", fs))
// }


// ServeFiles sets up static file serving
func ServeFiles() {
    rootDir, err := filepath.Abs("../")
    if err != nil {
        log.Printf("Error getting root directory: %v", err)
        return
    }
    
    staticPath := filepath.Join(rootDir, "static")
    log.Printf("Static files directory: %s", staticPath)
    
    // Verify static directory exists
    if _, err := os.Stat(staticPath); os.IsNotExist(err) {
        log.Printf("Static directory not found: %s", staticPath)
        return
    }

    // Verify CSS file exists
    cssPath := filepath.Join(staticPath, "css", "style.css")
    if _, err := os.Stat(cssPath); os.IsNotExist(err) {
        log.Printf("CSS file not found: %s", cssPath)
    } else {
        log.Printf("CSS file found: %s", cssPath)
    }
    
    fs := http.FileServer(http.Dir(staticPath))
    
    // Wrap file server with logging
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Static request received: %s", r.URL.Path)
        requestedFile := filepath.Join(staticPath, r.URL.Path)
        log.Printf("Looking for file: %s", requestedFile)
        
        if _, err := os.Stat(requestedFile); os.IsNotExist(err) {
            log.Printf("File not found: %s", requestedFile)
            http.NotFound(w, r)
            return
        }
        
        log.Printf("Serving file: %s", requestedFile)
        fs.ServeHTTP(w, r)
    })
    
    http.Handle("/static/", http.StripPrefix("/static/", handler))
}