package middleware

import (
		"log"
    "net/http"
)

// SecurityHeaders adds security-related HTTP headers to all responses
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Define security headers
			headers := map[string]string{
				"X-Frame-Options":           "SAMEORIGIN",
				"X-XSS-Protection":          "1; mode=block",
				"X-Content-Type-Options":    "nosniff",
				"Referrer-Policy":           "strict-origin-when-cross-origin",
				"Content-Security-Policy":   "default-src 'self' https:; " +
																		"script-src 'self' 'unsafe-inline' https://unpkg.com; " +
																		"style-src 'self' 'unsafe-inline' https:; " +
																		"img-src 'self' https: data:; " +
																		"frame-src 'self' https://www.youtube.com; " +
																		"connect-src 'self' https:;",
				"Permissions-Policy":        "microphone=(), camera=()",  // Removed unnecessary restrictions
				"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		}

        // Set all security headers
        for key, value := range headers {
            w.Header().Set(key, value)
        }

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