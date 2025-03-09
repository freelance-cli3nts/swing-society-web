// internal/routes.go
package internal

import (
    "log"
		"os"
		"encoding/json"
    "net/http"
    "path/filepath"
    "strings"
    "swing-society-website/server/internal/api/handlers"
    "swing-society-website/server/internal/config"
    "swing-society-website/server/internal/middleware"
    "swing-society-website/server/internal/storage"
)

type Router struct {
    mw          *middleware.MiddlewareManager
    rateLimiter *middleware.IPRateLimiter
}

func NewRouter(projectID string) (*Router, error) {
    r := &Router{
        mw:          middleware.NewMiddlewareManager(),
        rateLimiter: middleware.NewIPRateLimiter(),
    }

    // Add global middleware
    r.mw.UseGlobal(
        middleware.SecurityHeaders,
        middleware.Logger,
    )

    return r, nil
}

func createFileServer(dir, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the requested path
			path := r.URL.Path
			fullPath := filepath.Join(dir, path)
			
			log.Printf("Base static directory: %s", dir)
			log.Printf("Requested path: %s", path)
			
			// Determine content type based on file extension if not specified
			if contentType == "" {
					ext := filepath.Ext(path)
					switch ext {
					case ".css":
							contentType = "text/css"
					case ".js":
							contentType = "application/javascript"
					case ".html":
							contentType = "text/html"
					case ".json":
							contentType = "application/json"
					}
			}
			
			// Set content type if determined
			if contentType != "" {
					w.Header().Set("Content-Type", contentType)
			}
			
			// Check if file exists
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					http.NotFound(w, r)
					return
			}
			
			// Serve the file
			http.ServeFile(w, r, fullPath)
	})
}


func DirectCSSHandler(staticDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
			// Extract the CSS file path
			path := strings.TrimPrefix(r.URL.Path, "/css/")
			fullPath := filepath.Join(staticDir, "css", path)
			
			// Log the path for debugging
			log.Printf("Serving CSS file: %s from %s", path, fullPath)
			
			// Always set the correct MIME type
			w.Header().Set("Content-Type", "text/css")
			
			// Check if file exists
			_, err := os.Stat(fullPath)
			if os.IsNotExist(err) {
					log.Printf("CSS file not found: %s", fullPath)
					http.NotFound(w, r)
					return
			}
			
			// Read the file content
			content, err := os.ReadFile(fullPath)
			if err != nil {
					log.Printf("Error reading CSS file: %v", err)
					http.Error(w, "Error reading CSS file", http.StatusInternalServerError)
					return
			}
			
			// Write the CSS content directly
			w.Write(content)
	}
}


func (r *Router) SetupRoutes() error {
    // Create middleware chains
    apiChain := middleware.NewChain().
        Add(r.rateLimiter.RateLimitMiddleware).
				Add(middleware.TimestampMiddleware)

    staticChain := middleware.NewChain().
        Add(r.rateLimiter.RateLimitMiddleware)

    // Initialize all storage implementations
    googleStorage := storage.NewGoogleSheetsStorage(config.AppConfig.External.GoogleSheetsURL)
    jsonStorage := storage.NewJSONFileStorage(filepath.Join(config.AppConfig.Paths.DataDir, "carousel.json"))
    registrationStorage := storage.NewSimpleRegistrationStorage()
    classStorage := storage.NewSimpleClassStorage()
    templateStorage := storage.NewFileTemplateStorage()
    newsletterStorage := storage.NewSimpleNewsletterStorage()
    contactStorage := storage.NewSimpleContactStorage()
    
    // Initialize all handlers with their dependencies
    carouselHandler := handlers.NewCarouselHandler(googleStorage, jsonStorage)
    registrationHandler := handlers.NewRegistrationHandler(registrationStorage)
    classHandler := handlers.NewClassHandler(classStorage)
    templateHandler := handlers.NewTemplateHandler(templateStorage)
    newsletterHandler := handlers.NewNewsletterHandler(newsletterStorage)
    contactHandler := handlers.NewContactHandler(contactStorage)

    // Health check endpoint (no rate limiting)
    r.mw.AddHandler("/health", http.HandlerFunc(handlers.HandleHealth))

		// Firebase health check endpoint
		r.mw.AddHandler("/health/firebase", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Force Firebase check by adding query param
		r.URL.Query().Set("check_firebase", "true")
		handlers.HandleHealth(w, r)
		}))


		r.mw.AddHandler("/api/test-firebase", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Initialize Firebase client
			firebase, err := storage.NewFirebaseClient()
			if err != nil {
					http.Error(w, "Failed to initialize Firebase: "+err.Error(), http.StatusInternalServerError)
					return
			}
			
			// Test connection
			data, err := firebase.TestConnection(r.Context())
			if err != nil {
					http.Error(w, "Failed to connect to Firebase: "+err.Error(), http.StatusInternalServerError)
					return
			}
			
			// Write success response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
					"status": "success",
					"message": "Firebase connection established",
					"data": data,
			})
	}), apiChain.GetMiddlewares()...)



    // API routes (with rate limiting)
		
		// Routes for dynamic email validation
		r.mw.AddHandler("/api/validate-email", 
				http.HandlerFunc(registrationHandler.ValidateEmail), 
				apiChain.GetMiddlewares()...)
		// Routes for dynamic phone validation
		r.mw.AddHandler("/api/validate-phone", 
				http.HandlerFunc(registrationHandler.ValidatePhone), 
				apiChain.GetMiddlewares()...)
		// Routes for dynamic name validation
		r.mw.AddHandler("/api/validate-name", 
				http.HandlerFunc(registrationHandler.ValidateName), 
				apiChain.GetMiddlewares()...)

    r.mw.AddHandler("/api/carousel/", 
        http.HandlerFunc(carouselHandler.ServeCarousel), 
        apiChain.GetMiddlewares()...)
    
    r.mw.AddHandler("/api/register", 
        http.HandlerFunc(registrationHandler.HandleRegistration), 
        apiChain.GetMiddlewares()...)
    
    r.mw.AddHandler("/api/class", 
        http.HandlerFunc(classHandler.HandleClass), 
        apiChain.GetMiddlewares()...)
    
    r.mw.AddHandler("/api/newsletter", 
        http.HandlerFunc(newsletterHandler.HandleNewsletter), 
        apiChain.GetMiddlewares()...)
    
    r.mw.AddHandler("/api/contact", 
        http.HandlerFunc(contactHandler.HandleContact), 
        apiChain.GetMiddlewares()...)

		// In your SetupRoutes method, replace the static handlers with:
		staticDir := config.AppConfig.Paths.StaticDir

		// Static file handlers (with caching and rate limiting)
		staticHandler := http.StripPrefix("/static/", createFileServer(staticDir, ""))
		r.mw.AddHandler("/static/", staticHandler, staticChain.GetMiddlewares()...)

		// CSS files handler with proper MIME type
		cssHandler := http.StripPrefix("/css/", DirectCSSHandler(staticDir))
		r.mw.AddHandler("/css/", cssHandler, staticChain.GetMiddlewares()...)

		// JavaScript files handler with proper MIME type
		jsHandler := http.StripPrefix("/js/", createFileServer(staticDir, "application/javascript"))
		r.mw.AddHandler("/js/", jsHandler, staticChain.GetMiddlewares()...)

    // Template route
    r.mw.AddHandler("/templates/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        templateHandler.HandleTemplate(w, r, r.URL.Path)
    }))

    // Root handler - must be registered last to avoid overriding other routes
    r.mw.AddHandler("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip API and other specific paths
        if strings.HasPrefix(r.URL.Path, "/api/") || 
           strings.HasPrefix(r.URL.Path, "/static/") ||
           strings.HasPrefix(r.URL.Path, "/css/") ||
           strings.HasPrefix(r.URL.Path, "/js/") ||
           strings.HasPrefix(r.URL.Path, "/templates/") {
            return
        }
        
        log.Printf("Root handler serving index.html for path: %s", r.URL.Path)
        indexPath := filepath.Join(config.AppConfig.Paths.TemplatesDir, "index.html")
        http.ServeFile(w, r, indexPath)
    }))

    // Register all handlers with http package
    registeredHandlers := r.mw.GetHandlers()
    for _, h := range registeredHandlers {
        handler := r.mw.BuildHandler(h.Path, h.Handler, h.Middlewares...)
        http.Handle(h.Path, handler)
    }


    return nil
}
