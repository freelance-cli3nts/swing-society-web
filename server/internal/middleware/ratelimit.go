package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"context"
	"swing-society-website/server/internal/config"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	ips map[string]*visitor
	mu  *sync.RWMutex
}

func NewIPRateLimiter() *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*visitor),
		mu:  &sync.RWMutex{},
	}

	go i.cleanupVisitors()
	return i
}

// parseRateLimit parses rate limit strings like "100-M" (100 per minute)
func parseRateLimit(limit string) (rate.Limit, int) {
	parts := strings.Split(limit, "-")
	if len(parts) != 2 {
		return rate.Limit(1), 1 // default fallback
	}

	requestsPerPeriod, _ := strconv.Atoi(parts[0])
	period := strings.ToUpper(parts[1])

	switch period {
	case "S":
		return rate.Limit(requestsPerPeriod), requestsPerPeriod
	case "M":
		return rate.Every(time.Minute / time.Duration(requestsPerPeriod)), requestsPerPeriod
	case "H":
		return rate.Every(time.Hour / time.Duration(requestsPerPeriod)), requestsPerPeriod
	default:
		return rate.Limit(1), 1
	}
}

func (i *IPRateLimiter) GetLimiter(ip string, path string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	key := ip + path
	v, exists := i.ips[key]

	if !exists {
		// Get rate limit from config based on path
		limitStr := config.AppConfig.Security.RateLimits["default"]

		// Apply specific limits based on path
		if strings.HasPrefix(path, "/api/") {
			limitStr = config.AppConfig.Security.RateLimits["api"]
		} else if strings.HasPrefix(path, "/static/") {
			limitStr = config.AppConfig.Security.RateLimits["static"]
		}

		rateLimit, burst := parseRateLimit(limitStr)
		limiter := rate.NewLimiter(rateLimit, burst)

		v = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		i.ips[key] = v
		log.Printf("New rate limiter created for IP: %s, Path: %s", ip, path)
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (i *IPRateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute * 5)
		i.mu.Lock()
		for key, v := range i.ips {
			if time.Since(v.lastSeen) > time.Hour {
				delete(i.ips, key)
			}
		}
		i.mu.Unlock()
	}
}

// RateLimitMiddleware is an HTTP middleware that enforces rate limiting.
// If the request exceeds the allowed rate, it responds with HTTP 429.
func (i *IPRateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			limiter := i.GetLimiter(ip, r.URL.Path)
			if !limiter.Allow() {
					http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
					return
			}
			next.ServeHTTP(w, r)
	})
}


// TimestampMiddleware adds the current timestamp to the request context
func TimestampMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add current time to context
			ctx := context.WithValue(r.Context(), "requestTime", time.Now().Unix())
			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
	})
}