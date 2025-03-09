// internal/middleware/chain.go
package middleware

import (
    "log"
    "net/http"
    "time"
)

// Middleware represents a middleware function
type Middleware func(http.Handler) http.Handler

// Chain holds a list of middleware to be executed in order
type Chain struct {
    middlewares []Middleware
}

// NewChain creates a new middleware chain
func NewChain() *Chain {
    return &Chain{
        middlewares: make([]Middleware, 0),
    }
}

// Add appends a middleware to the chain
func (c *Chain) Add(middleware Middleware) *Chain {
    c.middlewares = append(c.middlewares, middleware)
    return c
}

// GetMiddlewares returns the slice of middlewares
func (c *Chain) GetMiddlewares() []Middleware {
    return c.middlewares
}

// AddIf adds a middleware conditionally
func (c *Chain) AddIf(condition bool, middleware Middleware) *Chain {
    if condition {
        return c.Add(middleware)
    }
    return c
}

// Then executes the middleware chain and ends with the final handler
func (c *Chain) Then(h http.Handler) http.Handler {
    if h == nil {
        h = http.DefaultServeMux
    }

    // Execute middlewares in reverse order
    for i := len(c.middlewares) - 1; i >= 0; i-- {
        h = c.middlewares[i](h)
    }

    return h
}

// Handler represents any HTTP handler with additional context
type Handler struct {
    Path        string
    Handler     http.Handler
    Middlewares []Middleware
}

// MiddlewareManager manages middleware chains for different routes
type MiddlewareManager struct {
    globalChain *Chain
    handlers    []Handler
}

// NewMiddlewareManager creates a new middleware manager
func NewMiddlewareManager() *MiddlewareManager {
    return &MiddlewareManager{
        globalChain: NewChain(),
        handlers:    make([]Handler, 0),
    }
}

// UseGlobal adds middleware to the global chain
func (m *MiddlewareManager) UseGlobal(middleware ...Middleware) {
    for _, mw := range middleware {
        m.globalChain.Add(mw)
    }
}

// AddHandler adds a handler with its specific middleware
func (m *MiddlewareManager) AddHandler(path string, handler http.Handler, middleware ...Middleware) {
    m.handlers = append(m.handlers, Handler{
        Path:        path,
        Handler:     handler,
        Middlewares: middleware,
    })
}

// BuildHandler builds the final handler with all middleware applied
func (m *MiddlewareManager) BuildHandler(path string, handler http.Handler, middleware ...Middleware) http.Handler {
    // Create a new chain with global middleware
    chain := NewChain()
    
    // Add logging middleware
    chain.Add(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            log.Printf("Started %s %s", r.Method, r.URL.Path)
            next.ServeHTTP(w, r)
            log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
        })
    })

    // Add global middleware
    for _, mw := range m.globalChain.middlewares {
        chain.Add(mw)
    }

    // Add route-specific middleware
    for _, mw := range middleware {
        chain.Add(mw)
    }

    return chain.Then(handler)
}

func (m *MiddlewareManager) GetHandlers() []Handler {
	return m.handlers
}

