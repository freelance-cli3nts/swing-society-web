package internal

import (
    "log"
    "net/http"
    "strings"
    "path/filepath"
    "os"
    "swinng-society-website/server/internal/config"
)

type Router struct {
    // ... keeping minimal required fields
}

func NewRouter(projectID string) (*Router, error) {
    return &Router{}, nil
}

func (r *Router) SetupRoutes() {
    // Serve CSS files
    http.HandleFunc("/css/", func(w http.ResponseWriter, req *http.Request) {
        file := filepath.Join(config.AppPaths.StaticDir, req.URL.Path)
        w.Header().Set("Content-Type", "text/css")
        http.ServeFile(w, req, file)
    })

    // Serve JavaScript files
    http.HandleFunc("/js/", func(w http.ResponseWriter, req *http.Request) {
        file := filepath.Join(config.AppPaths.StaticDir, req.URL.Path)
        w.Header().Set("Content-Type", "application/javascript")
        http.ServeFile(w, req, file)
    })

    // Handle static assets
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.AppPaths.StaticDir))))

    // Handle template requests - THIS IS THE IMPORTANT PART
    http.HandleFunc("/templates/", func(w http.ResponseWriter, req *http.Request) {
        log.Printf("Template request received: %s", req.URL.Path)
        log.Printf("Is HTMX request: %v", req.Header.Get("HX-Request"))

        templatePath := filepath.Join(config.AppPaths.RootDir, req.URL.Path)
        
        // Check if file exists
        if _, err := os.Stat(templatePath); os.IsNotExist(err) {
            log.Printf("Template not found: %s", templatePath)
            http.NotFound(w, req)
            return
        }

        // Handle HTMX request
        if req.Header.Get("HX-Request") == "true" {
            log.Printf("Serving template content: %s", templatePath)
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            content, err := os.ReadFile(templatePath)
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }
            w.Write(content)
            return
        }

        // For non-HTMX requests, serve index.html
        indexPath := filepath.Join(config.AppPaths.TemplatesDir, "index.html")
        http.ServeFile(w, req, indexPath)
    })

    // Handle root and other paths
    http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        // Ignore API routes
        if strings.HasPrefix(req.URL.Path, "/api/") {
            return
        }

        // Serve index.html for all other routes
        indexPath := filepath.Join(config.AppPaths.TemplatesDir, "index.html")
        log.Printf("Serving index.html for path: %s", req.URL.Path)
        http.ServeFile(w, req, indexPath)
    })
}