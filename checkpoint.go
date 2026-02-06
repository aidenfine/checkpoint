package checkpoint

import (
	"net/http"
)

func Limit(maxTokens, refillRate, tokensPerRefill int, config Config) func(next http.Handler) http.Handler {
	return NewTokenBucket(maxTokens, refillRate, tokensPerRefill, config).Handler
}

func LimitByIp(maxTokens, refillRate, tokensPerRefill int, config Config) func(next http.Handler) http.Handler {
	return Limit(maxTokens, refillRate, tokensPerRefill, config)
}
func LimitIpByEndpoint(maxTokens, refillRate, tokensPerRefill int, config Config) func(next http.Handler) http.Handler {
	return Limit(maxTokens, refillRate, tokensPerRefill, config)
}

func WithConfig(config Config) func(next http.Handler) http.Handler {
	return config.LimitMethod(config.MaxTokens, config.RefillRate, config.TokensPerRefill, config)
}
