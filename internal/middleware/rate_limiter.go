package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests per period
	period   time.Duration // time period
}

// Visitor represents a rate limit visitor
type Visitor struct {
	lastSeen  time.Time
	requests  []time.Time
	mu        sync.Mutex
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed
// period: time window (e.g., 1 minute)
func NewRateLimiter(rate int, period time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		period:   period,
	}

	// Cleanup old visitors every 5 minutes
	go rl.cleanupVisitors()

	return rl
}

// getVisitor gets or creates a visitor
func (rl *RateLimiter) getVisitor(ip string) *Visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &Visitor{
			lastSeen: time.Now(),
			requests: make([]time.Time, 0),
		}
		rl.visitors[ip] = v
	}

	return v
}

// allow checks if a request should be allowed
func (rl *RateLimiter) allow(ip string) bool {
	visitor := rl.getVisitor(ip)
	visitor.mu.Lock()
	defer visitor.mu.Unlock()

	now := time.Now()
	visitor.lastSeen = now

	// Remove requests outside the time window
	cutoff := now.Add(-rl.period)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range visitor.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	visitor.requests = validRequests

	// Check if limit exceeded
	if len(visitor.requests) >= rl.rate {
		return false
	}

	// Add current request
	visitor.requests = append(visitor.requests, now)
	return true
}

// cleanupVisitors removes old visitors
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			v.mu.Lock()
			if time.Since(v.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
			v.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware returns a middleware that rate limits requests
func RateLimitMiddleware(rate int, period time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, period)

	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		// Check rate limit
		if !limiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Rate limit exceeded. Please try again later.",
				"code":    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// StrictRateLimitMiddleware for sensitive endpoints (login, register, password reset)
func StrictRateLimitMiddleware() gin.HandlerFunc {
	// 5 requests per minute for sensitive endpoints
	return RateLimitMiddleware(5, 1*time.Minute)
}

// ModerateRateLimitMiddleware for normal API endpoints
func ModerateRateLimitMiddleware() gin.HandlerFunc {
	// 60 requests per minute for normal endpoints
	return RateLimitMiddleware(60, 1*time.Minute)
}
