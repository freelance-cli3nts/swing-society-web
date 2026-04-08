package internal

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "net/http"
    "path/filepath"
    "strings"
    "swing-society-website/server/internal/api/handlers"
    "swing-society-website/server/internal/config"
    "swing-society-website/server/internal/email"
    "swing-society-website/server/internal/gcal"
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
					case ".svg":
							contentType = "image/svg+xml"
					}
			}
			
			// Set content type if determined
			if contentType != "" {
					w.Header().Set("Content-Type", contentType)
			}
			
			log.Printf("Serving %s with Content-Type: %s", path, contentType)

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

func SVGHandler(staticDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
			// Extract the SVG file path
			path := strings.TrimPrefix(r.URL.Path, "/images/")
			fullPath := filepath.Join(staticDir, "assets", "images", path)
			
			// Log the path for debugging
			log.Printf("Serving SVG file: %s from %s", path, fullPath)
			
			ext := filepath.Ext(path)
			if ext == ".svg" {
					w.Header().Set("Content-Type", "image/svg+xml")
					log.Printf("Serving SVG file: %s with MIME type image/svg+xml", path)
			}

			// Check if file exists
			_, err := os.Stat(fullPath)
			if os.IsNotExist(err) {
					log.Printf("SVG file not found: %s", fullPath)
					http.NotFound(w, r)
					return
			}
			
			// Read the file content
			content, err := os.ReadFile(fullPath)
			if err != nil {
					log.Printf("Error reading SVG file: %v", err)
					http.Error(w, "Error reading SVG file", http.StatusInternalServerError)
					return
			}
			
			// Write the SVG content directly
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

    // Initialize email service
    emailSvc := email.NewService(config.AppConfig)

    // Initialize Google Calendar service (optional — degrades gracefully if not configured)
    gcalSvc, gcalErr := gcal.NewService(config.AppConfig.External.CalendarID)
    if gcalErr != nil {
        log.Printf("Warning: Google Calendar init failed: %v", gcalErr)
    }

    // Initialize all storage implementations
    googleStorage := storage.NewGoogleSheetsStorage(config.AppConfig.External.GoogleSheetsURL)
    jsonStorage := storage.NewJSONFileStorage(filepath.Join(config.AppConfig.Paths.DataDir, "carousel.json"))
    classStorage := storage.NewSimpleClassStorage()
    templateStorage := storage.NewFileTemplateStorage()

    firebase, err := storage.NewFirebaseClient()
    if err != nil {
        return fmt.Errorf("Firebase initialization failed: %v", err)
    }
    registrationStorage := storage.NewSimpleRegistrationStorage(firebase)
    newsletterStorage := storage.NewSimpleNewsletterStorage(firebase)
    contactStorage := storage.NewSimpleContactStorage(firebase)
    eventNotificationStorage := storage.NewSimpleEventNotificationStorage(firebase)
    calendarStorage := storage.NewCalendarStorage(firebase)

    // Initialize all handlers with their dependencies
    carouselHandler := handlers.NewCarouselHandler(googleStorage, jsonStorage)
    registrationHandler := handlers.NewRegistrationHandler(registrationStorage, emailSvc)
    classHandler := handlers.NewClassHandler(classStorage)
    templateHandler := handlers.NewTemplateHandler(templateStorage)
    newsletterHandler := handlers.NewNewsletterHandler(newsletterStorage, emailSvc)
    contactHandler := handlers.NewContactHandler(contactStorage, emailSvc)
    eventNotificationHandler := handlers.NewEventNotificationHandler(eventNotificationStorage)
    calendarHandler := handlers.NewCalendarHandler(calendarStorage, gcalSvc)

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
        
    // Test endpoint for Firebase schema
    r.mw.AddHandler("/api/test-firebase-schema", 
        http.HandlerFunc(storage.TestFirebaseSchemaHandler), 
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
        
    r.mw.AddHandler("/api/notifications",
        http.HandlerFunc(eventNotificationHandler.HandleEventNotification),
        apiChain.GetMiddlewares()...)

    r.mw.AddHandler("/api/calendar/",
        http.HandlerFunc(calendarHandler.HandleCalendar),
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

		// Images handler with proper MIME type
		svgHandler := http.StripPrefix("/images/", SVGHandler(staticDir))
		r.mw.AddHandler("/images/", svgHandler, staticChain.GetMiddlewares()...)

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
