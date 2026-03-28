package rate

import (
	"context"
	"testing"
	"time"
)

func TestLocalLimiterEnforcesWindow(t *testing.T) {
	limiter := NewLocalLimiter()
	ctx := context.Background()

	allowed, remaining, _, err := limiter.Allow(ctx, "user:1", 2, time.Minute)
	if err != nil {
		t.Fatalf("first allow failed: %v", err)
	}
	if !allowed || remaining != 1 {
		t.Fatalf("unexpected first result: allowed=%v remaining=%d", allowed, remaining)
	}

	allowed, remaining, _, err = limiter.Allow(ctx, "user:1", 2, time.Minute)
	if err != nil {
		t.Fatalf("second allow failed: %v", err)
	}
	if !allowed || remaining != 0 {
		t.Fatalf("unexpected second result: allowed=%v remaining=%d", allowed, remaining)
	}

	allowed, remaining, _, err = limiter.Allow(ctx, "user:1", 2, time.Minute)
	if err != nil {
		t.Fatalf("third allow failed: %v", err)
	}
	if allowed || remaining != 0 {
		t.Fatalf("expected limiter to block on third request: allowed=%v remaining=%d", allowed, remaining)
	}
}
