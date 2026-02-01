package checkpoint_test

import (
	"testing"
	"time"

	checkpoint "github.com/aidenfine/checkpoint"
)

func TestTokenBucketLimiter_HasEnoughTokens(t *testing.T) {

	limiter := checkpoint.NewTokenBucket(100, 1, 5)

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

	limiter := checkpoint.NewTokenBucket(100, 1, 5)

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

	limiter := checkpoint.NewTokenBucket(100, 1, 5)

	ip := "1"

	limiter.SetClientForTest(ip, 0, time.Now().AddDate(0, 0, -1))
	allowed, _ := limiter.Allow(ip)

	if !allowed {
		t.Errorf("expected request to be allowed but has failed")
	}
}
