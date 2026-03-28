package proxy

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portlyn/internal/domain"
	"portlyn/internal/routing"
)

func TestProxyPassThroughAndFailureModes(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "upstream-ok")
	}))
	defer upstream.Close()

	manager := NewManager(newFakeRoutingStore(routing.RouteConfig{
		ID:          "1",
		ServiceID:   1,
		ServiceName: "svc",
		Host:        "app.example.com",
		Path:        "/",
		TargetURL:   upstream.URL,
		Service: domain.Service{
			ID:   1,
			Name: "svc",
			Domain: domain.Domain{
				Name: "app.example.com",
			},
		},
		EffectivePolicy: domain.AccessPolicy{AccessMode: domain.AccessModePublic},
	}), NewInMemoryConfigCache(), NewInMemoryConfigBus(), nil, nil, nil, nil, ManagerOptions{
		LocalCacheTTL:      time.Hour,
		LocalCacheCapacity: 16,
	})

	req := httptest.NewRequest(http.MethodGet, "http://app.example.com/", nil)
	req.Host = "app.example.com"
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	manager.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected proxy success, got %d: %s", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != "upstream-ok" {
		t.Fatalf("unexpected proxy body %q", rec.Body.String())
	}

	downManager := NewManager(newFakeRoutingStore(routing.RouteConfig{
		ID:          "2",
		ServiceID:   2,
		ServiceName: "down",
		Host:        "down.example.com",
		Path:        "/",
		TargetURL:   "http://127.0.0.1:1",
		Service: domain.Service{
			ID:   2,
			Name: "down",
			Domain: domain.Domain{
				Name: "down.example.com",
			},
		},
		EffectivePolicy: domain.AccessPolicy{AccessMode: domain.AccessModePublic},
	}), NewInMemoryConfigCache(), NewInMemoryConfigBus(), nil, nil, nil, nil, ManagerOptions{
		LocalCacheTTL:      time.Hour,
		LocalCacheCapacity: 16,
	})

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://down.example.com/", nil)
		req.Host = "down.example.com"
		req.RemoteAddr = "127.0.0.1:12345"
		rec := httptest.NewRecorder()
		downManager.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusBadGateway {
			t.Fatalf("expected bad gateway during upstream failures, got %d", rec.Code)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "http://down.example.com/", nil)
	req.Host = "down.example.com"
	req.RemoteAddr = "127.0.0.1:12345"
	rec = httptest.NewRecorder()
	downManager.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected degraded 503 after repeated upstream failures, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestManagerInvalidateHostPublishesConfigUpdate(t *testing.T) {
	store := newFakeRoutingStore(testRouteConfig("1", "example.com", "/"))
	manager := NewManager(store, NewInMemoryConfigCache(), NewInMemoryConfigBus(), nil, nil, nil, nil, ManagerOptions{
		LocalCacheTTL:      time.Hour,
		LocalCacheCapacity: 16,
	})
	if err := manager.InvalidateHost(context.Background(), "example.com"); err != nil {
		t.Fatalf("invalidate host: %v", err)
	}
}
