package checkpointmiddleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type ClientRequestData struct {
	LastRequest time.Time
	Tokens      int
}

type TokenBucket struct {
	mu              sync.Mutex
	clients         map[string]ClientRequestData
	tokensPerRefill int
	refillRate      int // seconds per token
	maxTokens       int
	onRateLimited   http.HandlerFunc
}

func NewTokenBucket(refillRate, maxTokens int, tokensPerRefill int) *TokenBucket {
	return &TokenBucket{
		clients:         make(map[string]ClientRequestData),
		refillRate:      refillRate,
		maxTokens:       maxTokens,
		tokensPerRefill: tokensPerRefill,
	}
}

func (tb *TokenBucket) getClient(ip string) ClientRequestData {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.clients[ip]
}

func (tb *TokenBucket) setClient(ip string, data ClientRequestData) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.clients[ip] = data
}

// TODO: is there a better/cleaner way of doing this? Seems like a weird way of doing testing
func (tb *TokenBucket) SetClientForTest(ip string, tokens int, lastRequest time.Time) {
	tb.setClient(ip, ClientRequestData{
		LastRequest: lastRequest,
		Tokens:      tokens,
	})
}

func (tb *TokenBucket) Allow(ip string) (bool, int) {
	client := tb.getClient(ip)
	now := time.Now()

	if client.LastRequest.IsZero() {
		tb.setClient(ip, ClientRequestData{LastRequest: now, Tokens: tb.maxTokens - 1})
		return true, tb.maxTokens - 1
	}

	elapsed := now.Sub(client.LastRequest).Seconds()
	tokensToAdd := (int(elapsed) / tb.refillRate) * tb.tokensPerRefill
	newTokens := min(client.Tokens+tokensToAdd, tb.maxTokens)

	if newTokens <= 0 {
		return false, 0
	}

	newTokens--
	tb.setClient(ip, ClientRequestData{LastRequest: now, Tokens: newTokens})
	return true, newTokens
}

func (tb *TokenBucket) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		allowed, remainingTokens := tb.Allow(ip)

		if !allowed {
			if tb.onRateLimited != nil {
				tb.onRateLimited(w, r)
				return
			}
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.Header().Set("X-RateLimit-Remaining", string(rune(remainingTokens)))

		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the real client IP from the request
// This handles cases where the app is behind a reverse proxy
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by most reverse proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs (client, proxy1, proxy2, ...)
		// The first one is the original client
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header (used by nginx and others)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fallback to RemoteAddr
	// This will be the proxy's IP if behind a reverse proxy without proper headers
	ip := r.RemoteAddr
	// RemoteAddr includes port, so we need to strip it
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return ip
}
