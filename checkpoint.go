package checkpoint

import (
	"errors"
	"net/http"
	"time"
)

// entry point of the go application
// user can call LimitByIp and then it handles everything for them

// add options here
type Option func(rl *RateLimiter)
type GetKey func(r *http.Request) (string, error)
type GetKeyWithStringArg func(r *http.Request, s string) (string, error)

// build limiter
func Limit(requestLimit int, window time.Duration, opts ...Option) func(h http.Handler) http.Handler {
	rl := CreateRateLimiter(requestLimit, window, opts...)
	return rl.Handler
}

func LimitByIp(requestLimit int, window time.Duration) func(h http.Handler) http.Handler {
	return Limit(requestLimit, window, WithKey(IpKeyFunc))
}
func LimitByEndpoint(requestLimit int, window time.Duration) func(h http.Handler) http.Handler {
	return Limit(requestLimit, window, WithKey(EndpointKeyFunc))
}

func WithKey(keyFunc GetKey) Option {
	return func(rl *RateLimiter) {
		rl.keyFunc = keyFunc
	}
}
func IpKeyFunc(r *http.Request) (string, error) {
	return r.RemoteAddr, nil
}
func EndpointKeyFunc(r *http.Request) (string, error) {
	return r.URL.Path, nil
}
func UserIdKeyFunc(r *http.Request, userId string) (string, error) {
	if userId != "" {
		return userId, nil
	}
	return "", errors.New("empty userId")
}
