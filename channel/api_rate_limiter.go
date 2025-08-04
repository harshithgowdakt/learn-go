package main

import (
    "net/http"
    "time"
)

type RateLimiter struct {
    tokens   chan struct{}
    refill   *time.Ticker
    capacity int
}

func NewRateLimiter(capacity int, refillRate time.Duration) *RateLimiter {
    rl := &RateLimiter{
        tokens:   make(chan struct{}, capacity),
        refill:   time.NewTicker(refillRate),
        capacity: capacity,
    }
    
    // Fill bucket initially
    for i := 0; i < capacity; i++ {
        rl.tokens <- struct{}{}
    }
    
    // Start refill goroutine
    go rl.refillTokens()
    
    return rl
}

func (rl *RateLimiter) refillTokens() {
    for range rl.refill.C {
        select {
        case rl.tokens <- struct{}{}:
            // Token added
        default:
            // Bucket full, skip
        }
    }
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}

func (rl *RateLimiter) Wait() {
    <-rl.tokens
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !rl.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Usage
func main() {
    limiter := NewRateLimiter(10, 100*time.Millisecond) // 10 requests per second
    
    http.Handle("/api", limiter.Middleware(http.HandlerFunc(apiHandler)))
    http.ListenAndServe(":8080", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("API response"))
}