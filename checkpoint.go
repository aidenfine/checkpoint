package checkpoint

import (
	"net/http"
)

func Limit(maxTokens, refillRate, tokensPerRefill int) func(next http.Handler) http.Handler {
	return NewTokenBucket(maxTokens, refillRate, tokensPerRefill).Handler
}

func LimitByIp(maxTokens, refillRate, tokensPerRefill int) func(next http.Handler) http.Handler {
	return Limit(maxTokens, refillRate, tokensPerRefill)
}
func LimitIpByEndpoint(maxTokens, refillRate, tokensPerRefill int) func(next http.Handler) http.Handler {
	return Limit(maxTokens, refillRate, tokensPerRefill)
}

func WithConfig(config Config) func(next http.Handler) http.Handler {
	return config.LimitMethod(config.MaxTokens, config.RefillRate, config.TokensPerRefill)
}
