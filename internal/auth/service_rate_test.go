package auth

import (
	"context"
	"testing"
	"time"

	"portlyn/internal/config"
)

type fakeRateLimiter struct {
	allowed bool
	calls   int
}

func (l *fakeRateLimiter) Allow(context.Context, string, int, time.Duration) (bool, int, time.Time, error) {
	l.calls++
	return l.allowed, 0, time.Now().UTC().Add(time.Minute), nil
}

func TestServiceUsesConfiguredRateLimiter(t *testing.T) {
	service, err := NewService(nil, nil, nil, nil, nil, "12345678901234567890123456789012", "portlyn", "", time.Minute, time.Minute, config.OIDCConfig{}, config.OTPConfig{}, time.Minute, config.RateLimitConfig{
		LoginAttempts: 1,
		Window:        time.Minute,
	}, time.Minute, true, nil)
	if err != nil {
		t.Fatalf("new service failed: %v", err)
	}

	limiter := &fakeRateLimiter{allowed: false}
	service.SetRateLimiter(limiter)

	if !service.isRateLimited(context.Background(), "user:1", time.Now().UTC()) {
		t.Fatalf("expected configured limiter to block request")
	}
	if limiter.calls != 1 {
		t.Fatalf("expected limiter to be called once, got %d", limiter.calls)
	}
}
