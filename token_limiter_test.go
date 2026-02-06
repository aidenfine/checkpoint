package checkpoint_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/aidenfine/checkpoint"
)

func generateRandomIPv4() string {
	o1 := rand.Intn(256)
	o2 := rand.Intn(256)
	o3 := rand.Intn(256)
	o4 := rand.Intn(256)

	return fmt.Sprintf("%d.%d.%d.%d", o1, o2, o3, o4)
}

func TestTokenBucketLimiter_HasEnoughTokens(t *testing.T) {

	config := checkpoint.Config{
		IgnorePaths:     []string{},
		MaxTokens:       100,
		RefillRate:      1,
		TokensPerRefill: 5,
		LimitMethod:     checkpoint.LimitByIp,
	}

	limiter := checkpoint.NewTokenBucket(100, 1, 5, config)

	ip := "1"

	allowed, remaining := limiter.Allow(ip)

	if !allowed {
		t.Errorf("expected request to be allowed but failed")
	}
	if remaining != 99 {
		t.Errorf("expected remaining to be 99 but got %d", remaining)
	}
}
func TestTokenBucketLimiter_DoesNotHaveEnoughTokens(t *testing.T) {

	config := checkpoint.Config{
		IgnorePaths:     []string{},
		MaxTokens:       100,
		RefillRate:      1,
		TokensPerRefill: 5,
		LimitMethod:     checkpoint.LimitByIp,
	}

	limiter := checkpoint.NewTokenBucket(100, 1, 5, config)

	ip := "1"

	limiter.SetClientForTest("1", 0, time.Now())

	allowed, remaining := limiter.Allow(ip)

	if allowed {
		t.Errorf("expected request to fail but was allowed, has %d requests remaining", remaining)
	}
	if remaining != 0 {
		t.Errorf("expected remaining to be 0 but got %d", remaining)
	}
}

func TestTokenBucketLimiter_HasZeroTokensButCanBeRefilled(t *testing.T) {

	config := checkpoint.Config{
		IgnorePaths:     []string{},
		MaxTokens:       100,
		RefillRate:      1,
		TokensPerRefill: 5,
		LimitMethod:     checkpoint.LimitByIp,
	}

	limiter := checkpoint.NewTokenBucket(100, 1, 5, config)

	ip := "1"

	limiter.SetClientForTest(ip, 0, time.Now().AddDate(0, 0, -1))
	allowed, _ := limiter.Allow(ip)

	if !allowed {
		t.Errorf("expected request to be allowed but has failed")
	}
}

// ---------------- BENCHMARKS ----------------//
func BenchmarkAllowUniqueIps(b *testing.B) {

	config := checkpoint.Config{
		IgnorePaths:     []string{},
		MaxTokens:       100,
		RefillRate:      1,
		TokensPerRefill: 1,
		LimitMethod:     checkpoint.LimitByIp,
	}
	tb := checkpoint.NewTokenBucket(100, 1, 1, config)

	ips := make([]string, 1000)
	for i := range ips {
		ips[i] = generateRandomIPv4()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.Allow(ips[i%len(ips)])
	}
}
func BenchmarkAllowNonUniqueIps(b *testing.B) {
	config := checkpoint.Config{
		IgnorePaths:     []string{},
		MaxTokens:       100,
		RefillRate:      1,
		TokensPerRefill: 1,
		LimitMethod:     checkpoint.LimitByIp,
	}
	tb := checkpoint.NewTokenBucket(100, 1, 1, config)

	ips := make([]string, 5)
	for i := range ips {
		ips[i] = generateRandomIPv4()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.Allow(ips[i%len(ips)])
	}
}
