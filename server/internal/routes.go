package internal

import (
    "log"
    "net/http"
    "path/filepath"
    "strings"
    "swing-society-website/server/internal/api/handlers"
    "swing-society-website/server/internal/config"
)

type Router struct {}

func NewRouter(projectID string) (*Router, error) {
    return &Router{}, nil
}

func (r *Router) SetupRoutes() {
    // Static file handlers
    http.HandleFunc("/css/", serveStatic("text/css"))
    http.HandleFunc("/js/", serveStatic("application/javascript"))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.AppPaths.StaticDir))))

    // API routes
    http.HandleFunc("/api/register", handlers.HandleRegistration)
    http.HandleFunc("/api/contact", handlers.HandleContact)
    http.HandleFunc("/api/newsletter", handlers.HandleNewsletter)
    http.HandleFunc("/api/class", handlers.HandleClass)

    // Template routes
    http.HandleFunc("/templates/", handleTemplateRequest)

    // Root handler
    http.HandleFunc("/", handleRoot)
}

func serveStatic(contentType string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        file := filepath.Join(config.AppPaths.StaticDir, r.URL.Path)
        w.Header().Set("Content-Type", contentType)
        http.ServeFile(w, r, file)
    }
}

func handleTemplateRequest(w http.ResponseWriter, r *http.Request) {
    log.Printf("Template request received: %s", r.URL.Path)
    
    if r.Header.Get("HX-Request") == "true" {
        handlers.HandleTemplate(w, r, r.URL.Path)
        return
    }

    // Serve index.html for non-HTMX requests
    indexPath := filepath.Join(config.AppPaths.TemplatesDir, "index.html")
    http.ServeFile(w, r, indexPath)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    if strings.HasPrefix(r.URL.Path, "/api/") {
        return // Let API handlers handle these routes
    }

    indexPath := filepath.Join(config.AppPaths.TemplatesDir, "index.html")
    log.Printf("Serving index.html for path: %s", r.URL.Path)
    http.ServeFile(w, r, indexPath)
}


// package internal

// import (
//     "log"
//     "net/http"
//     "strings"
//     "path/filepath"
//     "os"
//     "swinng-society-website/server/internal/config"
// 		"swinng-society-website/server/internal/handlers/forms"
// )

// type Router struct {
//     // ... keeping minimal required fields
// }

// func NewRouter(projectID string) (*Router, error) {
//     return &Router{}, nil
// }

// func (r *Router) SetupRoutes() {
//     // Serve CSS files
//     http.HandleFunc("/css/", func(w http.ResponseWriter, req *http.Request) {
//         file := filepath.Join(config.AppPaths.StaticDir, req.URL.Path)
//         w.Header().Set("Content-Type", "text/css")
//         http.ServeFile(w, req, file)
//     })

//     // Serve JavaScript files
//     http.HandleFunc("/js/", func(w http.ResponseWriter, req *http.Request) {
//         file := filepath.Join(config.AppPaths.StaticDir, req.URL.Path)
//         w.Header().Set("Content-Type", "application/javascript")
//         http.ServeFile(w, req, file)
//     })

//     // Handle static assets
//     http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.AppPaths.StaticDir))))

//     // Handle template requests - THIS IS THE IMPORTANT PART
//     http.HandleFunc("/templates/", func(w http.ResponseWriter, req *http.Request) {
//         log.Printf("Template request received: %s", req.URL.Path)
//         log.Printf("Is HTMX request: %v", req.Header.Get("HX-Request"))

//         templatePath := filepath.Join(config.AppPaths.RootDir, req.URL.Path)
        
//         // Check if file exists
//         if _, err := os.Stat(templatePath); os.IsNotExist(err) {
//             log.Printf("Template not found: %s", templatePath)
//             http.NotFound(w, req)
//             return
//         }

//         // Handle HTMX request
//         if req.Header.Get("HX-Request") == "true" {
//             log.Printf("Serving template content: %s", templatePath)
//             w.Header().Set("Content-Type", "text/html; charset=utf-8")
//             content, err := os.ReadFile(templatePath)
//             if err != nil {
//                 http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//                 return
//             }
//             w.Write(content)
//             return
//         }

//         // For non-HTMX requests, serve index.html
//         indexPath := filepath.Join(config.AppPaths.TemplatesDir, "index.html")
//         http.ServeFile(w, req, indexPath)
//     })

//     // Handle root and other paths
//     http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
//         // Ignore API routes
//         if strings.HasPrefix(req.URL.Path, "/api/") {
//             return
//         }

//         // Serve index.html for all other routes
//         indexPath := filepath.Join(config.AppPaths.TemplatesDir, "index.html")
//         log.Printf("Serving index.html for path: %s", req.URL.Path)
//         http.ServeFile(w, req, indexPath)
//     })

// 		// Form handling routes
// 		formHandler := handlers.NewFormHandler()
// 		http.HandleFunc("/register", formHandler.HandleRegistration)

// 		// Handle API routes
// 		http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
// 				// Add your API route handling here
// 				w.Header().Set("Content-Type", "application/json")
// 				// ... rest of your API handling
// 		})
// }