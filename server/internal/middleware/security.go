package middleware

import (
		"log"
    "net/http"
		"strings"
		"swing-society-website/server/internal/config"
)

// SecurityHeaders adds security-related HTTP headers to all responses
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Use CSP directives from config
        w.Header().Set("Content-Security-Policy", config.AppConfig.Security.CSPDirectives)
        
        // Use allowed origins from config
        if len(config.AppConfig.Security.AllowedOrigins) > 0 {
            w.Header().Set("Access-Control-Allow-Origin", 
                strings.Join(config.AppConfig.Security.AllowedOrigins, ", "))
        }

        // Standard security headers
        w.Header().Set("X-Frame-Options", "SAMEORIGIN")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        w.Header().Set("Permissions-Policy", "microphone=(), camera=()")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

        next.ServeHTTP(w, r)
    })
}

// Logger middleware for request logging
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Log request details
        log.Printf("[%s] %s %s %s", 
            r.RemoteAddr, 
            r.Method, 
            r.URL.Path,
            r.UserAgent(),
        )
        
        next.ServeHTTP(w, r)
    })
}