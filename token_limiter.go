package checkpoint

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
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
	ignorePaths     []string
}

func NewTokenBucket(maxTokens, refillRate int, tokensPerRefill int, config Config) *TokenBucket {
	tb := &TokenBucket{
		clients:         make(map[string]ClientRequestData),
		refillRate:      refillRate,
		maxTokens:       maxTokens,
		tokensPerRefill: tokensPerRefill,
		ignorePaths:     config.IgnorePaths,
	}
	return tb
}

func (tb *TokenBucket) SetClientForTest(ip string, tokens int, lastRequest time.Time) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.clients[ip] = ClientRequestData{
		LastRequest: lastRequest,
		Tokens:      tokens,
	}
}

func (tb *TokenBucket) Allow(ip string) (bool, int) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	client, exists := tb.clients[ip]
	now := time.Now()

	if !exists {
		tb.clients[ip] = ClientRequestData{
			LastRequest: now,
			Tokens:      tb.maxTokens - 1,
		}
		return true, tb.maxTokens - 1
	}

	elapsed := now.Sub(client.LastRequest).Seconds()
	tokensToAdd := (int(elapsed) / tb.refillRate) * tb.tokensPerRefill
	newTokens := min(client.Tokens+tokensToAdd, tb.maxTokens)

	if newTokens <= 0 {
		return false, 0
	}

	newTokens--
	tb.clients[ip] = ClientRequestData{
		LastRequest: now,
		Tokens:      newTokens,
	}

	return true, newTokens
}

func (tb *TokenBucket) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// is api path to be ignored? we may need to make this better because ex: logs/** w
		// ill not match logs/log we also need to ignore static assets like /favicon
		apiURL := r.URL.Path

		fmt.Printf("current URL: %s \n", apiURL)
		if slices.Contains(tb.ignorePaths, apiURL) {
			// directly serve http and do not limit
			next.ServeHTTP(w, r)
			return
		}

		ip, err := getClientIP(r)
		if err != nil {
			fmt.Println("error getting client ip:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		allowed, remainingTokens := tb.Allow(ip)

		if !allowed {
			if tb.onRateLimited != nil {
				tb.onRateLimited(w, r)
				return
			}
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remainingTokens))

		next.ServeHTTP(w, r)
	})
}
