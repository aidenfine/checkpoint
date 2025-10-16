package checkpoint_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aidenfine/checkpoint"
)

func TestRateLimiter_AllowsRequestUnderLimit(t *testing.T) {
	limiter := checkpoint.CreateRateLimiter(5, time.Second)

	handler := limiter.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Result().StatusCode)
	}
}

func TestRateLimiter_AllowRequestsFromDifferentIps(t *testing.T) {
	limiter := checkpoint.LimitByEndpoint(1, time.Second)

	handler := limiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	userReq := httptest.NewRequest("GET", "/users", nil)
	homeReq := httptest.NewRequest("GET", "/", nil)

	ur1 := httptest.NewRecorder()
	handler.ServeHTTP(ur1, userReq)

	ur2 := httptest.NewRecorder()
	handler.ServeHTTP(ur2, userReq)
	if ur2.Result().StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 Too many Requests, got %d", ur2.Result().StatusCode)
	}

	hr1 := httptest.NewRecorder()
	handler.ServeHTTP(hr1, homeReq)

	if hr1.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", hr1.Result().StatusCode)
	}
}
func TestRateLimiter_AllowRequestsFromDifferentSubRoutes(t *testing.T) {
	limiter := checkpoint.LimitByEndpoint(1, time.Second)

	handler := limiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	userReq := httptest.NewRequest("GET", "/users", nil)
	homeReq := httptest.NewRequest("GET", "/users/200", nil)

	ur1 := httptest.NewRecorder()
	handler.ServeHTTP(ur1, userReq)

	ur2 := httptest.NewRecorder()
	handler.ServeHTTP(ur2, userReq)
	if ur2.Result().StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 Too many Requests, got %d", ur2.Result().StatusCode)
	}

	hr1 := httptest.NewRecorder()
	handler.ServeHTTP(hr1, homeReq)

	if hr1.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", hr1.Result().StatusCode)
	}

}

func TestRateLimiter_BlocksRequestOverLimit(t *testing.T) {
	limiter := checkpoint.CreateRateLimiter(2, time.Second)

	handler := limiter.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)

	// First request
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req)

	// Second request
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req)

	// Expect this to be blocked
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, req)

	if w3.Result().StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 Too Many Requests, got %d", w3.Result().StatusCode)
	}
}

func TestRateLimiter_SetsRateLimitHeaders(t *testing.T) {
	limiter := checkpoint.CreateRateLimiter(3, time.Second, checkpoint.WithHeaders(checkpoint.ResponseHeaders{
		RateLimit:          "X-RateLimit-Limit",
		RateLimitRemaining: "X-RateLimit-Remaining",
	}))

	handler := limiter.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.Header.Get("X-RateLimit-Limit") != "3" {
		t.Errorf("expected X-RateLimit-Limit to be 3, got %s", resp.Header.Get("X-RateLimit-Limit"))
	}

	if resp.Header.Get("X-RateLimit-Remaining") != "2" {
		t.Errorf("expected X-RateLimit-Remaining to be 2, got %s", resp.Header.Get("X-RateLimit-Remaining"))
	}
}

func TestRateLimiter_CustomKeyFunc(t *testing.T) {
	limiter := checkpoint.CreateRateLimiter(1, time.Second, checkpoint.WithKey(func(r *http.Request) (string, error) {
		return "custom-user-id", nil
	}))

	handler := limiter.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)

	// First request should pass
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req)

	// Second request should be blocked
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req)

	if w2.Result().StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 Too Many Requests, got %d", w2.Result().StatusCode)
	}
}
func TestRateLimiter_RequestAfterWindow(t *testing.T) {
	limiter := checkpoint.CreateRateLimiter(1, 500*time.Millisecond)

	handler := limiter.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)

	// first request will pass
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req)
	if w1.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w1.Result().StatusCode)
	}

	// Second request will fail because of rate limit
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req)
	if w2.Result().StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w2.Result().StatusCode)
	}

	// Set time to let window expire
	time.Sleep(501 * time.Millisecond)

	// Third Request should be okay
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, req)
	if w3.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK after window, got %d", w3.Result().StatusCode)
	}
}
