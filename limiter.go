package checkpoint

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type LimitCount interface {
	Config(requestLimit int, window time.Duration)
	Get(key string, currWindow, prevWindow time.Time) (int, int, error)
}
type ResponseHeaders struct {
	Reset              string
	RateLimit          string // RateLimit-Limit
	RateLimitRemaining string // RateLimit-Remaining
	RateLimitReset     string // RateLimit-Reset
	RetryAfter         string
}

type Metadata struct {
	Ip  string
	URL string
}

type RateLimiter struct {
	requestLimit  int
	window        time.Duration
	limitCount    LimitCount
	onRateLimited http.HandlerFunc
	onError       func(http.ResponseWriter, *http.Request, error)
	headers       ResponseHeaders
	mu            sync.Mutex
	keyFunc       GetKey // define what method to rate limit user ex: ip, endpoint, userid...
}

type SlidingWindowLimitCount struct {
	mu      sync.Mutex
	clients map[string][]time.Time
	window  time.Duration
	limit   int
}

func CreateRateLimiter(requestLimit int, window time.Duration, opts ...Option) *RateLimiter {
	rl := &RateLimiter{
		requestLimit: requestLimit,
		window:       window,
		headers: ResponseHeaders{
			Reset:              "X-RateLimit-Reset",
			RateLimit:          "X-RateLimit-Limit",
			RateLimitRemaining: "X-RateLimit-Remaining",
			RateLimitReset:     "X-RateLimit-Reset",
			RetryAfter:         "X-RateLimit-Retry-After",
		},
	}
	for _, opt := range opts {
		opt(rl)
	}
	if rl.limitCount == nil {
		rl.limitCount = NewSlidingWindowLimitCount(requestLimit, window)
		rl.limitCount.Config(requestLimit, window)

	}
	if rl.onRateLimited == nil {
		rl.onRateLimited = rateLimitedResponse
	}
	if rl.onError == nil {
		rl.onError = onError
	}
	return rl
}

func NewSlidingWindowLimitCount(limit int, window time.Duration) *SlidingWindowLimitCount {
	return &SlidingWindowLimitCount{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (s *SlidingWindowLimitCount) Config(requestLimit int, window time.Duration) {
	s.limit = requestLimit
	s.window = window
}
func (s *SlidingWindowLimitCount) Get(key string, currWindow, prevWindow time.Time) (int, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-s.window)

	timestamps := s.clients[key]

	fmt.Println(timestamps, "timestamps")

	idx := 0
	for i, t := range timestamps {
		if t.After(windowStart) {
			idx = i
			break
		}
	}
	timestamps = timestamps[idx:]

	timestamps = append(timestamps, now)
	s.clients[key] = timestamps

	return len(timestamps), 0, nil
}
func (rl *RateLimiter) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key string
		var err error

		// default to ip
		if rl.keyFunc != nil {
			key, err = rl.keyFunc(r)
			if err != nil {
				rl.onError(w, r, err)
				return
			}
		} else {
			key = r.RemoteAddr // ip adddress
		}
		rl.mu.Lock()
		defer rl.mu.Unlock()
		count, _, err := rl.limitCount.Get(key, time.Now(), time.Now().Add(-rl.window))

		if err != nil {
			rl.onError(w, r, err)
			return
		}

		if count > rl.requestLimit {
			rl.onRateLimited(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// responses
func rateLimitedResponse(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
}

func onError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusPreconditionFailed)
}
