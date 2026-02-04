package checkpoint

import "net/http"

type Config struct {
	IgnorePaths     []string `json:"ignorePaths" yaml:"ignorePaths"`
	MaxTokens       int      `json:"maxTokens" yaml:"maxTokens"`
	RefillRate      int      `json:"refillRate" yaml:"refillRate"`
	TokensPerRefill int      `json:"tokensPerRefill" yaml:"tokensPerRefill"`
	LimitMethod     func(maxTokens, refillRate, tokensPerRefill int, config Config) func(next http.Handler) http.Handler
}
