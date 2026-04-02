package middleware

import (
		"fmt"
		"log"
    "net/http"
		"strings"
		"swing-society-website/server/internal/config"
)

// SecurityHeaders adds security-related HTTP headers to all responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Define the CSP directives with font-src included
			cspDirectives := "default-src 'self' https:; " +
					"script-src 'self' 'unsafe-inline' https://unpkg.com https://cdn.tailwindcss.com; " +
					"style-src 'self' 'unsafe-inline' https:; " +
					"img-src 'self' https: data:; " +
					"font-src 'self' data: https:; " +  // Added font-src directive
					"frame-src 'self' https://www.youtube.com; " +
					"connect-src 'self' https:;"
			
			// Use CORS settings from config
			if len(config.AppConfig.CORS.Origins) > 0 {
					w.Header().Set("Access-Control-Allow-Origin", 
							strings.Join(config.AppConfig.CORS.Origins, ", "))
			} else if len(config.AppConfig.Security.AllowedOrigins) > 0 {
					// Fallback to old config format for backward compatibility
					w.Header().Set("Access-Control-Allow-Origin", 
							strings.Join(config.AppConfig.Security.AllowedOrigins, ", "))
			}

			// Set CORS methods if configured
			if len(config.AppConfig.CORS.Methods) > 0 {
					w.Header().Set("Access-Control-Allow-Methods", 
							strings.Join(config.AppConfig.CORS.Methods, ", "))
			}

			// Set CORS headers if configured
			if len(config.AppConfig.CORS.ResponseHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", 
							strings.Join(config.AppConfig.CORS.ResponseHeaders, ", "))
			}

			// Set max age if configured
			if config.AppConfig.CORS.MaxAgeSeconds > 0 {
					w.Header().Set("Access-Control-Max-Age", 
							fmt.Sprintf("%d", config.AppConfig.CORS.MaxAgeSeconds))
			}

			// Standard security headers
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", cspDirectives)
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