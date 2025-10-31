package checkpoint

import (
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

// TODO: is there a more go way of doing this? Seems like a weird way of doing testing
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
