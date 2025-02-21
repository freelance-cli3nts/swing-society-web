package internal

import (
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
    handlers    map[string]http.Handler
    middleware  map[string]*middleware.Chain
}

func New(projectID string) (*Router, error) {
    r := &Router{
        mw:          middleware.NewMiddlewareManager(),
        rateLimiter: middleware.NewIPRateLimiter(),
        handlers:    make(map[string]http.Handler),
        middleware:  make(map[string]*middleware.Chain),
    }

    // Initialize middleware chains
    r.setupMiddlewareChains()
    
    // Initialize storage and handlers
    r.setupHandlers()

    return r, nil
}

func (r *Router) setupMiddlewareChains() {
    // Global middleware
    r.mw.UseGlobal(
        middleware.SecurityHeaders,
        middleware.Logger,
    )

    // API middleware chain
    r.middleware["api"] = middleware.NewChain().
        Add(r.rateLimiter.RateLimitMiddleware)

    // Static middleware chain
    r.middleware["static"] = middleware.NewChain().
        Add(r.rateLimiter.RateLimitMiddleware)
}

func (r *Router) setupHandlers() {
    // Initialize storages
    storages := r.initializeStorages()
    
    // Initialize handlers
    handlers := r.initializeHandlers(storages)
    
    // Store handlers for routes
    r.handlers = handlers
}

func (r *Router) initializeStorages() map[string]interface{} {
    return map[string]interface{}{
        "carousel":     []storage.CarouselStorage{
            storage.NewGoogleSheetsStorage(config.AppConfig.External.GoogleSheetsURL),
            storage.NewJSONFileStorage(filepath.Join(config.AppConfig.Paths.DataDir, "carousel.json")),
        },
        "registration": storage.NewSimpleRegistrationStorage(),
        "class":        storage.NewSimpleClassStorage(),
        "template":     storage.NewFileTemplateStorage(),
        "newsletter":   storage.NewSimpleNewsletterStorage(),
        "contact":      storage.NewSimpleContactStorage(),
    }
}

func (r *Router) initializeHandlers(storages map[string]interface{}) map[string]http.Handler {
    return map[string]http.Handler{
        "/health":          http.HandlerFunc(handlers.HandleHealth),
        "/api/carousel/":   http.HandlerFunc(handlers.NewCarouselHandler(
            storages["carousel"].([]storage.CarouselStorage)[0],
            storages["carousel"].([]storage.CarouselStorage)[1],
        ).ServeCarousel),
        "/api/register":    http.HandlerFunc(handlers.NewRegistrationHandler(
            storages["registration"].(storage.RegistrationStorage),
        ).HandleRegistration),
        "/api/class":      http.HandlerFunc(handlers.NewClassHandler(
            storages["class"].(storage.ClassStorage),
        ).HandleClass),
        "/api/newsletter": http.HandlerFunc(handlers.NewNewsletterHandler(
            storages["newsletter"].(storage.NewsletterStorage),
        ).HandleNewsletter),
        "/api/contact":    http.HandlerFunc(handlers.NewContactHandler(
            storages["contact"].(storage.ContactStorage),
        ).HandleContact),
				"/templates/": 		 http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    		handlers.NewTemplateHandler(storages["template"].(storage.TemplateStorage)).
        HandleTemplate(w, r, r.URL.Path)}),
        "/static/":        http.StripPrefix("/static/", 
            http.FileServer(http.Dir(config.AppConfig.Paths.StaticDir))),
        "/css/":          http.StripPrefix("/css/", 
            http.FileServer(http.Dir(config.AppConfig.Paths.StaticDir))),
        "/js/":           http.StripPrefix("/js/", 
            http.FileServer(http.Dir(config.AppConfig.Paths.StaticDir))),
    }
}

func (r *Router) SetupRoutes() error {
    // Register all routes with appropriate middleware
    for path, handler := range r.handlers {
        chain := r.getMiddlewareChain(path)
        r.mw.AddHandler(path, handler, chain.GetMiddlewares()...)
    }

    // Register all handlers with http package
    registeredHandlers := r.mw.GetHandlers()
    for _, h := range registeredHandlers {
        handler := r.mw.BuildHandler(h.Path, h.Handler, h.Middlewares...)
        http.Handle(h.Path, handler)
    }

    return nil
}

func (r *Router) getMiddlewareChain(path string) *middleware.Chain {
    if strings.HasPrefix(path, "/api/") {
        return r.middleware["api"]
    }
    if strings.HasPrefix(path, "/static/") || 
       strings.HasPrefix(path, "/css/") || 
       strings.HasPrefix(path, "/js/") {
        return r.middleware["static"]
    }
    return middleware.NewChain() // Empty chain for other routes
}