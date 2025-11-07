package checkpointmiddleware

import (
	"net/http"
)

func Limit(maxTokens, refillRate, tokensPerRefill int) func(next http.Handler) http.Handler {
	return NewTokenBucket(maxTokens, refillRate, tokensPerRefill).Handler
}

func LimitByIp(maxTokens, refillRate, tokensPerRefill int) func(next http.Handler) http.Handler {
	return Limit(maxTokens, refillRate, tokensPerRefill)
}
