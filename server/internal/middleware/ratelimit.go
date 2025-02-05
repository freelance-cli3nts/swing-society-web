package middleware

import (
    "golang.org/x/time/rate"
    "net/http"
    "sync"
    "strings"
    "log"
    "time"
)

type IPRateLimiter struct {
    ips map[string]*visitor
    mu  *sync.RWMutex
}

type visitor struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

func NewIPRateLimiter() *IPRateLimiter {
    i := &IPRateLimiter{
        ips: make(map[string]*visitor),
        mu:  &sync.RWMutex{},
    }
    
    // Start cleanup routine
    go i.cleanupVisitors()
    
    return i
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
    i.mu.Lock()
    defer i.mu.Unlock()

    v, exists := i.ips[ip]
    if !exists {
        // More generous limits for SPA:
        // - 30 requests per minute (covers initial page load + assets)
        // - Burst of 10 for concurrent asset loading
        limiter := rate.NewLimiter(rate.Every(200*time.Second), 5)
        v = &visitor{
            limiter:  limiter,
            lastSeen: time.Now(),
        }
        i.ips[ip] = v
        log.Printf("New rate limiter created for IP: %s", ip)
    }
    v.lastSeen = time.Now()
    return v.limiter
}

func (i *IPRateLimiter) RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Only apply rate limiting to specific paths
        if !shouldRateLimit(r.URL.Path) {
            next.ServeHTTP(w, r)
            return
        }

        ip := getIP(r)
        limiter := i.GetLimiter(ip)
        
        if !limiter.Allow() {
            log.Printf("Rate limit exceeded for IP: %s on path: %s", ip, r.URL.Path)
            w.Header().Set("Retry-After", "2")
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// shouldRateLimit determines if the path should be rate limited
func shouldRateLimit(path string) bool {
    // Rate limit API endpoints and form submissions
    if strings.HasPrefix(path, "/api/") || 
       strings.HasPrefix(path, "/submit/") ||
       strings.HasPrefix(path, "/contact/") {
        return true
    }
    
    // Don't rate limit static assets
    if strings.HasPrefix(path, "/static/") {
        return false
    }
    
    // Rate limit the initial page load (but with generous limits)
    return path == "/" || path == "/index.html"
}

func getIP(r *http.Request) string {
    // Get IP from various headers or remote addr
    ip := r.Header.Get("X-Real-IP")
    if ip == "" {
        ip = r.Header.Get("X-Forwarded-For")
    }
    if ip == "" {
        ip = strings.Split(r.RemoteAddr, ":")[0]
    }
    return ip
}

func (i *IPRateLimiter) cleanupVisitors() {
    for {
        time.Sleep(5 * time.Minute)
        i.mu.Lock()
        for ip, v := range i.ips {
            if time.Since(v.lastSeen) > 10*time.Minute {
                delete(i.ips, ip)
                log.Printf("Cleaned up rate limiter for IP: %s", ip)
            }
        }
        i.mu.Unlock()
    }
}